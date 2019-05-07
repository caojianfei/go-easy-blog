package response

import (
	"encoding/json"
	"go-easy-blog/errors"
	"io"
	"net/http"
)

type Response interface {
	String() (string, error)
}

type SuccessResponse struct {
	Code int `json:"code"`
	Msg  string `json:"msg"`
	Data interface{} `json:"data"`
}

type FailResponse struct {
	Code int `json:"code"`
	Msg  string `json:"msg"`
}

func (res SuccessResponse) String() (string, error) {
	js, err := json.Marshal(res)
	return string(js), err
}

func (res FailResponse) String() (string, error) {
	js, err := json.Marshal(res)
	return string(js), err
}

func Success(data interface{}, w http.ResponseWriter) error {
	res:= &SuccessResponse{0, "success", data}
	str, err := res.String()
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = io.WriteString(w, str)
	return err
}

func Fail(code int, msg string, w http.ResponseWriter) error {
	res := &FailResponse{code, msg}
	str, err := res.String()

	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = io.WriteString(w, str)
	return err
}

// 业务上的系统错误
func SystemError(w http.ResponseWriter) error {
	return Fail(errors.SystemError, "系统错误", w)
}

// 业务上参数错误
func ParamError(w http.ResponseWriter, errMsg string) error {
	return Fail(errors.ParamError, errMsg, w)
}

// 通用业务错误
func BusinessError(w http.ResponseWriter, errMsg string) error {
	return Fail(errors.BusinessError, errMsg, w)
}

