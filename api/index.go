package handler

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

// List 会返回给交付层一个列表回应
type Result struct {
	Code    int         `json:"code"`    //请求状态代码
	Message interface{} `json:"message"` //请求结果提示
	Data    interface{} `json:"data"`    //请求结果
}

func Handler(w http.ResponseWriter, r *http.Request) {

	method := r.Method

	switch method {
	case "GET":
		checkHandler(w, r)
	case "POST":
		msgHandler(w, r)
	default:
		defaultHandler(w, r)
	}
}

func checkHandler(w http.ResponseWriter, r *http.Request) {
	// TODO 微信校验
	params := r.URL.Query()
	value := params.Get("echostr")
	w.Write([]byte(value))
}

func msgHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)

	msg := &WxReqMsg{}

	if err := xml.Unmarshal(body, &msg); err == nil {
		w.Write([]byte("success"))
	}
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("method not suooprt"))
}

type WxReqMsg struct {
	ToUserName   string `xml:"ToUserName"`
	FromUserName string `xml:"FromUserName"`
	CreateTime   int64  `xml:"CreateTime"`
	MsgType      string `xml:"MsgType"`
	Content      string `xml:"Content"`
	MsgId        string `xml:"MsgId"`
	MsgDataId    string `xml:"MsgDataId"`
	Idx          string `xml:"Idx"`
}
