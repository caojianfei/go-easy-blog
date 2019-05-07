package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/siddontang/go-log/log"
	"go-easy-blog/database"
	"go-easy-blog/response"
	"net/http"
	"strconv"
	"time"
)

type CreateTagData struct {
	Name string `json:"name"`
	Description string `json:"description"`
}

type TagItem struct {
	Id int `json:"id"`
	Name string `json:"name"`
	ArticleNumber int `json:"article_number"`
	Description string `json:"description"`
}

// 创建标签
func CreateTag(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	data := &CreateTagData{}
	err := json.NewDecoder(r.Body).Decode(data)
	// fmt.Printf("request data %v", data)
	if err != nil {
		log.Errorf("parse request body error, err: ", err.Error())
		response.SystemError(w)
		return
	}

	if data.Name == "" {
		log.Error("params error, no tag name")
		response.ParamError(w, "标签名称不能为空")
		return
	}

	// 查询标签是否存在
	db := database.New()
	stmp, err := db.Prepare(database.QueryTagByName)
	if err != nil {
		log.Errorf("prepare sql error. sql: %s, error: %s", database.QueryTagByName, err.Error())
		response.SystemError(w)
		return
	}

	defer stmp.Close()

	existTag := &database.TagModel{}
	err = stmp.QueryRow(data.Name).Scan(existTag.Id, existTag.Name, existTag.ArticleNumber, existTag.Description, existTag.CreatedAt, existTag.UpdatedAt)
	if err != nil{
		if err != sql.ErrNoRows {
			log.Errorf("query sql error. sql: %s, error %s", database.QueryTagByName, err.Error())
			response.SystemError(w)
			return
		}
	}

	if err == nil {
		log.Errorf("已经存在【%s】标签，请勿重复添加", data.Name)
		response.BusinessError(w, fmt.Sprintf("标签【%s】已经存在，请勿重复添加", data.Name))
		return
	}

	// 创建标签

	stmp, err = db.Prepare(database.InsertTag)
	if err != nil {
		log.Errorf("prepare sql error. sql: %s, error: %s", database.InsertTag)
		response.SystemError(w)
		return
	}
	defer stmp.Close()

	nowTime := time.Now().Format("2006-01-02 15:04:05")
	res, err := stmp.Exec(data.Name, data.Description, nowTime, nowTime)
	if err != nil {
		log.Errorf("exec sql error. sql: %s, error: %s", database.InsertTag, err.Error())
		response.SystemError(w)
		return
	}

	lastInsertId, err := res.LastInsertId()
	if err != nil {
		log.Errorf("get last insert id error: %s", err.Error())
	}

	response.Success(map[string]interface{}{"id": lastInsertId}, w)
	return;
}

// 删除标签
func DeleteTag(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	tagId := params.ByName("id")
	db := database.New()

	// 查询标签是否存在
	stmp, err := db.Prepare(database.QueryTagById)
	if err != nil {
		log.Errorf("prepare sql error. sql: %s, error: %s", database.QueryTagById, err.Error())
		response.SystemError(w)
		return
	}

	defer stmp.Close()

	existTag := database.TagModel{}
	err = stmp.QueryRow(tagId).Scan(&existTag.Id, &existTag.Name, &existTag.ArticleNumber, &existTag.Description, &existTag.CreatedAt, &existTag.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			response.BusinessError(w, "该标签不存在或已经被删除")
		} else {
			log.Errorf("query tag error: %s", err.Error())
			response.SystemError(w)
		}
		return
	}

	stmp, err = db.Prepare(database.SoftDeleteTagById)
	if err != nil {
		log.Errorf("prepare sql error. sql: %s, error: %s", database.SoftDeleteTagById, err.Error())
		response.SystemError(w)
		return
	}
	defer stmp.Close()

	nowTime := time.Now().Format("2006-01-02 15:04:05")

	_, err = stmp.Exec(nowTime, tagId)
	if err != nil {
		log.Errorf("delete tag fail. sql: %s, tagId: %s, error: %s", database.SoftDeleteTagById, tagId, err.Error())
		response.BusinessError(w, "标签删除失败")
		return
	}

	return
}

// 获取表现列表（分页）
func TagList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// fmt.Printf("query: %v \n", r.URL.Query())
	// 标签名搜索
	queryName, ok := r.URL.Query()["name"]
	var filterName string
	if ok {
		filterName = queryName[0]
	}
	// 每页数量
	var pageSize = 5
	queryPageSize, ok := r.URL.Query()["pageSize"]
	if ok {
		pageSizeStr := queryPageSize[0]
		if pageSizeInt, err := strconv.Atoi(pageSizeStr); err == nil {
			pageSize = pageSizeInt
		}
	}
	// 分页
	var page = 0
	queryPage, ok := r.URL.Query()["page"]
	if ok {
		pageStr := queryPage[0]
		if pageInt, err := strconv.Atoi(pageStr); err == nil {
			page = pageInt
		}
	}
	querySql := "SELECT `id`, `name`, `article_number`, `description` FROM tags "
	if len(filterName) > 0 {
		querySql += "WHERE name LIKE ? AND id > ? LIMIT ?"
	} else {
		querySql += "WHERE id > ? LIMIT ?"
	}

	// 获取列表
	db := database.New()
	stmp, err := db.Prepare(querySql)
	if err != nil {
		log.Errorf("prepare sql error. sql: %s, error: %s", querySql, err.Error())
		response.SystemError(w)
		return
	}
	defer stmp.Close()

	var rows *sql.Rows
	if len(filterName) > 0 {
		rows, err = stmp.Query("%" + filterName + "%", page, pageSize)
	} else {
		rows, err = stmp.Query(page, pageSize)
	}

	if err != nil {
		if err == sql.ErrNoRows {

		} else {
			log.Errorf("query sql error. sql: %s, error: %s", querySql, err.Error())
			response.SystemError(w)
			return
		}
	}
	defer rows.Close()

	var tagList = make([]TagItem, pageSize)
	i := 0
	for rows.Next() {
		err = rows.Scan(&tagList[i].Id, &tagList[i].Name, &tagList[i].ArticleNumber, &tagList[i].Description)
		if err != nil {
			log.Printf("scan tag error. error: %s", err.Error())
			response.SystemError(w)
			return
		}
		i++
	}

	if err = rows.Err(); err != nil {
		log.Print("query result error. error: %s", err.Error())
		response.SystemError(w)
		return
	}

	result := make(map[string]interface{})

	// 空集
	if i == 0 {
		result["list"] = tagList[:i]
		result["page"] = page
		result["end"] = true
	} else {
		result["list"] = tagList[:i]
		result["end"] = i < pageSize
		result["page"] = tagList[i - 1].Id
	}

	response.Success(result, w)
	// fmt.Printf("querySql: %s, filterName: %s, pageSize: %d, page: %d", querySql, filterName, pageSize, page)
}