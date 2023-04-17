package handler

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

const openapi_url = "https://api.openai.com/v1/chat/completions"

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

	rspMsg := &WxRspMsg{
		FromUserName: msg.ToUserName,
		ToUserName:   msg.FromUserName,
		CreateTime:   time.Now().Unix(),
		MsgType:      "text",
		Content:      "hello",
	}

	rspMsg.Content = SendMsg(msg.Content)

	result, err := xml.Marshal(rspMsg)
	if err != nil {
		w.Write([]byte("success"))
	}

	w.Write(result)

}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("method not suooprt"))
}

func SendMsg(msg string) string {

	reqMsg := &ChatReqMsg{
		Model: "gpt-3.5-turbo",
	}

	chatMsg := &ChatMsg{
		Role:    "user",
		Content: msg,
	}

	reqMsg.Messages = append(reqMsg.Messages, *chatMsg)

	data, err := json.Marshal(reqMsg)
	if err != nil {
		// TODO
		fmt.Printf("序列化错误 err=%v\n", err)
		return ""
	}

	req, _ := http.NewRequest("POST", openapi_url, strings.NewReader(string(data)))

	req.Header.Set("Authorization", "Bearer "+os.Getenv("API_KEY"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	respByte, _ := ioutil.ReadAll(resp.Body)

	var rspMsg ChatRspMsg

	// TODO 错误信息的数据结构跟正常不一样
	json.Unmarshal(respByte, &rspMsg)

	return rspMsg.Choices[0].Message.Content
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

type ChatReqMsg struct {
	Model    string    `json:"model"`
	Messages []ChatMsg `json:"messages"`
}

type ChatMsg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRspMsg struct {
	Id string `json:"id"`
	// Object string `json:object`
	Created int           `json:"created"`
	Model   string        `json:"model"`
	Choices []ChatChoices `json:"choices"`
}

type ChatChoices struct {
	Message      ChatMsg `json:"message"`
	FinishReason string  `json:"finish_reason"`
	Index        int32   `json:"index"`
}
