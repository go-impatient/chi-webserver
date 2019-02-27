package model

import (
	"net/http"

	"github.com/moocss/chi-webserver/src/pkg/errno"
	"github.com/moocss/chi-webserver/src/pkg/render"
)

// Result represents HTTP response body.
type Result struct {
	Code 			int         `json:"code"` // return code, 0 for succ
	Message  	string      `json:"message"`  // message
	Data 			interface{} `json:"data"` // data object
}

// 输出返回结果
func SendResult(w http.ResponseWriter, data interface{}, err error) {
	code, message := errno.DecodeErr(err)
	render.JSON(w, &Result{
		Code:    code,
		Message: message,
		Data:    data,
	}, http.StatusOK)
}
