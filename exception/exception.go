package exception

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// NewAPIException 创建一个API异常
// 用于其他模块自定义异常
func NewAPIException(namespace string, code int, reason, format string, a ...interface{}) *APIException {
	// 0表示正常状态, 但是要排除变量的零值
	if code == 0 {
		code = -1
	}

	var httpCode int
	if code/100 >= 1 && code/100 <= 5 {
		//1<=code<=5  100-500
		httpCode = code
	} else {
		httpCode = http.StatusInternalServerError
	}

	return &APIException{
		Namespace: namespace,
		ErrCode:   code,
		Reason:    reason,
		HttpCode:  httpCode,
		Message:   fmt.Sprintf(format, a...),
	}
}

// NewAPIException 创建一个API异常
// 用于其他模块自定义异常
func NewAPIException1(code int, reason, format string, a ...interface{}) *APIException {
	// 0表示正常状态, 但是要排除变量的零值
	if code == 0 {
		code = -1
	}

	var httpCode int
	if code/100 >= 1 && code/100 <= 5 {
		httpCode = code
	} else {
		httpCode = http.StatusInternalServerError
	}

	return &APIException{
		ErrCode:  code,
		Reason:   reason,
		HttpCode: httpCode,
		Message:  fmt.Sprintf(format, a...),
	}
}

// APIException API异常
type APIException struct {
	Namespace string `json:"namespace"`
	HttpCode  int    `json:"http_code"`
	ErrCode   int    `json:"error_code"`
	Reason    string `json:"reason"`
	Message   string `json:"message"`
	Meta      any    `json:"meta"`
	Data      any    `json:"data"`
}

// 调用这个方法返回结构体的信息
func (e *APIException) ToJson() string {
	dj, _ := json.Marshal(e)
	return string(dj)
}

// 返回错误信息
func (e *APIException) Error() string {
	return e.Message
}

// Code exception's code, 如果code不存在返回-1
func (e *APIException) ErrorCode() int {
	return int(e.ErrCode)
}

func (e *APIException) GetReason() string {
	return e.Reason
}

func (e *APIException) GetData() interface{} {
	return e.Data
}

func (e *APIException) GetMeta() interface{} {
	return e.Meta
}

// Code exception's code, 如果code不存在返回-1
func (e *APIException) GetHttpCode() int {
	return int(e.HttpCode)
}

func (e *APIException) GetNamespace() string {
	return e.Namespace
}

func (e *APIException) WithNamespace(ns string) *APIException {
	e.Namespace = ns
	return e
}
