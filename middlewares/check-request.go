package middlewares

import (
	"bytes"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/siddontang/go-log/log"
	"io"
	"io/ioutil"
	"net/http"
)

var AcceptContentType = "application/json"

// 检查请求消息体是否是json数据
func CheckRequest(next httprouter.Handle) httprouter.Handle {

	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
		// fmt.Println("CheckRequest Middleware")
		// 消息体为空不判断，get,head,delete 不判断
		if r.ContentLength == 0 || r.Method == http.MethodGet || r.Method == http.MethodHead || r.Method == http.MethodDelete {
			next(w, r, param)
			return
		}

		headers := r.Header
		if contentTypeArr, ok := headers["Content-Type"]; ok {
			//fmt.Println(contentTypeArr)
			contentType := contentTypeArr[0]
			// 不接受的请求体数据类型
			if contentType != AcceptContentType {
				checkRequestFail(w)
				return
			}
		} else {
			// 必须设置 Content-Type
			checkRequestFail(w)
			return
		}

		requestContent, err := ioutil.ReadAll(r.Body)
		_ = r.Body.Close()
		r.Body = ioutil.NopCloser(bytes.NewBuffer(requestContent))

		if err != nil {
			log.Error("read request body error")
		}

		//fmt.Printf("%v \n", requestContent)
		var contentJson interface{}
		err = json.Unmarshal(requestContent, &contentJson)
		//fmt.Printf("json content: %v\n", contentJson)
		if err != nil {
			log.Error(err)
			checkRequestFail(w)
			return
		}

		next(w, r, param)
	})
}

func checkRequestFail(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotAcceptable)
	_, _ = io.WriteString(w, " Not acceptable")
}