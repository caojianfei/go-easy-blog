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
	router.POST("/login",  middlewares.CheckRequest(handlers.Login))

	router.POST("/tag/create", middlewares.CheckLogin(middlewares.CheckRequest(handlers.CreateTag)))// 新增标签
	router.DELETE("/tag/delete/:id", middlewares.CheckLogin(handlers.DeleteTag))// 删除标签
	router.GET("/tag/list", middlewares.CheckLogin(handlers.TagList))// 标签列表

	router.POST("/article/create", middlewares.CheckLogin(middlewares.CheckRequest(handlers.CreateArticle)))
	router.GET("/article/:id", middlewares.CheckLogin(handlers.GetArticle))
	router.DELETE("/article/:id", middlewares.CheckLogin(handlers.DeleteArticle))
	return router
}

func TestServer(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Println("logined")
}
