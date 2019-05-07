package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/siddontang/go-log/log"
	"go-easy-blog/database"
	"go-easy-blog/helpers"
	"go-easy-blog/response"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type CreateArticleData struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Tags    []int  `json:"tags"`
	Status  int    `json:"status"`
}

// 新增文章
func CreateArticle(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	data := &CreateArticleData{}
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		log.Error("parse request body error. error: %s", err.Error())
		response.SystemError(w)
		return
	}

	// 获取认证token
	token, err := helpers.GetAuthorizationTokenFormRequestHeader(r)
	if err != nil {
		log.Error("get token error. error: %s", err.Error())
		response.SystemError(w)
		return
	}

	// 获取登录的用户id
	userId, err := helpers.GetUserIdFromAuthorization(token)
	if err != nil {
		log.Error("get userId error. error: %s", err.Error())
		response.SystemError(w)
		return
	}

	// 参数校验
	if data.Title == "" || data.Content == "" {
		log.Error("缺少参数")
		response.ParamError(w, "文章标题和内容不能为空")
		return
	}

	// 新增文章
	db := database.New()

	tx, err := db.Begin()
	if err != nil {
		log.Errorf("begin transaction error. error: %s", err.Error())
		response.SystemError(w)
		return
	}
	defer tx.Rollback()

	// 新增文章
	stmp, err := tx.Prepare(database.CreateNewArticle)
	if err != nil {
		log.Errorf("prepare sql error. sql: %s, error: %s", database.CreateNewArticle, err.Error())
		response.SystemError(w)
		return
	}

	nowTime := time.Now().Format("2006-01-02 15:04:05")
	res, err := stmp.Exec(userId, data.Title, data.Content, data.Status, nowTime, nowTime)
	if err != nil {
		log.Error("insert article error. sql: %s, error: %s", database.CreateNewArticle, err.Error())
		response.SystemError(w)
		return
	}

	lastInsertId, err := res.LastInsertId()
	if err != nil {
		log.Errorf("add article error. error: %s", err.Error())
		response.SystemError(w)
		return
	}

	_ = stmp.Close()

	// 文章-标签关联
	tags := data.Tags
	if len(tags) > 0 {
		insertSql := "INSERT INTO article_tag(article_id, tag_id) VALUES"
		values := make([]string, len(tags))
		for i := 0; i < len(tags); i++ {
			values[i] = fmt.Sprintf("(%d, %d)", lastInsertId, tags[i])
		}
		rows := strings.Join(values, ",")
		insertSql += rows
		_, err = db.Exec(insertSql)
		if err != nil {
			tx.Rollback()
			log.Errorf("prepare sql error. sql: %s, error: %s", insertSql, err.Error())
			response.SystemError(w)
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("commit transaction error. error: %s", err.Error())
		response.SystemError(w)
		return
	}

	result := map[string]int64{"id": lastInsertId}
	response.Success(result, w)
	return

}

type ArticleDetail struct {
	Id            int          `json:"id"`
	Title         string       `json:"title"`
	Content       string       `json:"content"`
	Status        int          `json:"status"`
	Views         int          `json:"views"`
	CommentNumber int          `json:"comment_number"`
	CreatedAt     string       `json:"created_at"`
	UpdatedAt     string       `json:"updated_at"`
	Author        ArticleUser  `json:"author"`
	Tags          []ArticleTag `json:"tags"`
}

type ArticleUser struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
}

type ArticleTag struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

// 获取文章详情
func GetArticle(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	articleId := params.ByName("id")
	articleIdInt, err := strconv.Atoi(articleId)
	if err != nil{
		response.ParamError(w, "参数错误")
		return
	}
	if articleIdInt < 1{
		response.ParamError(w, "参数错误")
		return
	}
	articleDetail := &ArticleDetail{}
	// 查询文章
	db := database.New()
	articleSql := "select a.`id` as `articleId`, a.`title`, a.`content`, a.`status`, a.`views`, a.`comment_number`, a.`created_at`, a.`updated_at`, b.`id` as `userId`, b.`username`, b.`nickname`" +
		" from articles as a inner join users as b on a.user_id=b.id where a.deleted_at is null and a.id = ?"
	stmp, err := db.Prepare(articleSql)
	if err != nil {
		log.Errorf("prepare sql error. sql: %s, error: %s", articleSql, err.Error())
		response.SystemError(w)
		return
	}
	defer stmp.Close()

	err = stmp.QueryRow(articleId).Scan(
		&articleDetail.Id,
		&articleDetail.Title,
		&articleDetail.Content,
		&articleDetail.Status,
		&articleDetail.Views,
		&articleDetail.CommentNumber,
		&articleDetail.CreatedAt,
		&articleDetail.UpdatedAt,
		&articleDetail.Author.Id,
		&articleDetail.Author.Username,
		&articleDetail.Author.Nickname,
	)
	if err != nil {
		log.Error("query article error. error: %s", err.Error())
		if err == sql.ErrNoRows {
			response.BusinessError(w, "文章不存在或已经被删除")
			return
		} else {
			response.SystemError(w)
			return
		}
	}

	// 获取文章标签
	tagSql := "select t.`id`, t.`name` from article_tag as a left join tags as t on a.tag_id=t.id where a.article_id = ?"
	stmp, err = db.Prepare(tagSql)
	if err != nil {
		log.Errorf("prepare sql error. sql: %s, error: %s", tagSql, err.Error())
		response.SystemError(w)
		return
	}
	defer stmp.Close()

	rows, err := stmp.Query(articleId)
	if err != nil {
		log.Errorf("query error. sql: %s, error: %s", tagSql, err.Error())
		response.SystemError(w)
		return
	}
	i := 0
	tags := make([]ArticleTag, 1)
	for rows.Next() {
		tag := ArticleTag{}
		err = rows.Scan(&tag.Id, &tag.Name)
		if err != nil {
			log.Errorf("scan tag error. error: %s", err.Error())
			response.SystemError(w)
			return
		}
		if i == 0 {
			tags[i] = tag
		} else {
			tags = append(tags, tag)
		}
		i++
	}
	fmt.Printf("tags: %v\n", tags)
	defer rows.Close()
	if err := rows.Err(); err != nil {
		log.Errorf("rows error. error: ", err.Error())
		response.SystemError(w)
		return
	}

	if i > 0 {
		articleDetail.Tags = tags
	}
	fmt.Printf("article %v \n", articleDetail)
	response.Success(articleDetail, w)
	return
}


// 删除文章
func DeleteArticle(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	articleId := params.ByName("id")
	articleIdInt, err := strconv.Atoi(articleId)
	if err != nil{
		response.ParamError(w, "参数错误")
		return
	}
	if articleIdInt < 1{
		response.ParamError(w, "参数错误")
		return
	}

	// 查询文章
	db := database.New()
	querySql := "select id from articles where id =? and deleted_at is null"
	stmp, err := db.Prepare(querySql)
	if err != nil {
		log.Errorf("prepare sql error. sql: %s, error: %s", querySql, err.Error())
		response.SystemError(w)
		return
	}
	defer stmp.Close()
	var id int
	err = stmp.QueryRow(articleId).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			response.BusinessError(w, "文章不存在或已经被删除")
			return
		} else {
			log.Errorf("query article error. sql: %s, error: %s", querySql, err.Error())
			response.SystemError(w)
			return
		}
	}

	updateSql := "update articles set deleted_at = ? where id = ?"
	nowTime := time.Now().Format("2006-01-02 15:04:05")

	stmp, err = db.Prepare(updateSql)
	if err != nil {
		log.Errorf("prepare sql error. sql: %s, error: %s", updateSql, err.Error())
		response.SystemError(w)
		return
	}

	_, err = stmp.Exec(nowTime, articleId)
	if err != nil {
		log.Error("delete article error. error: %s", err.Error())
		response.SystemError(w)
		return
	}

	response.Success(map[string]string{}, w)
	return
}