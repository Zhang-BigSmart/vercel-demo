package handler

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

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

	key := os.Getenv("API_KEY")

	rspMsg := &WxRspMsg{
		FromUserName: msg.ToUserName,
		ToUserName:   msg.FromUserName,
		CreateTime:   time.Now().Unix(),
		MsgType:      "text",
		Content:      "hello",
	}

	rspMsg.Content = key

	result, err := xml.Marshal(rspMsg)
	if err != nil {
		w.Write([]byte("success"))
	}

	w.Write(result)

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

type WxRspMsg struct {
	ToUserName   string   `xml:"ToUserName"`
	FromUserName string   `xml:"FromUserName"`
	CreateTime   int64    `xml:"CreateTime"`
	MsgType      string   `xml:"MsgType"`
	Content      string   `xml:"Content"`
	XMLName      xml.Name `xml:"xml"`
}
