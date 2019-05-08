package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/siddontang/go-log/log"
	"go-easy-blog/handlers"
	"go-easy-blog/middlewares"
	"net/http"
	"time"
)

func main() {

	now := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("now is: %s \n", now)

	// 注册路由
	router := RegisterRouter()

	// 启动 http 服务
	fmt.Println("http server start, port: 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func RegisterRouter() *httprouter.Router {
	router := httprouter.New()
	router.GET("/test", middlewares.CheckLogin(TestServer))

	// 管理员登录
	router.POST("/admin/login",  middlewares.CheckRequest(handlers.Login))

	// 标签管理
	router.POST("/admin/tag", middlewares.CheckLogin(middlewares.CheckRequest(handlers.CreateTag)))// 新增标签
	router.DELETE("/admin/tag/:id", middlewares.CheckLogin(handlers.DeleteTag))// 删除标签
	router.GET("/admin/tag", middlewares.CheckLogin(handlers.TagList))// 标签列表

	// 文章管理
	router.POST("/admin/article", middlewares.CheckLogin(middlewares.CheckRequest(handlers.CreateArticle)))// 新增文章
	router.POST("/admin/article/:id", middlewares.CheckLogin(middlewares.CheckRequest(handlers.UpdateArticle)))// 更新文章
	router.GET("/admin/article/:id", middlewares.CheckLogin(handlers.GetArticle))// 文章详情
	router.DELETE("/admin/article/:id", middlewares.CheckLogin(handlers.DeleteArticle)) // 删除文章


	// 访客路由


	return router
}

func TestServer(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Println("logined")
}
