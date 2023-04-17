package handler

import (
	"encoding/json"
	"net/http"
)

// List 会返回给交付层一个列表回应
type Result struct {
	Code    int         `json:"code"`    //请求状态代码
	Message interface{} `json:"message"` //请求结果提示
	Data    interface{} `json:"data"`    //请求结果
}

func Handler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json;charset=utf-8")

	data := "Hello! Smart Zhang!"

	response, _ := json.Marshal(&Result{
		Code:    200,
		Message: "ok",
		Data:    data,
	})

	_, _ = w.Write(response)

}
