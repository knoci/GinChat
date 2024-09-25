package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H struct {
	Code  int
	Msg   string
	Data  interface{}
	Rows  interface{}
	Total interface{}
}

func Resp(w http.ResponseWriter, code int, data interface{}, msg string) {
	// 设置响应的 Content-Type 头为 application/json，告诉客户端响应体是 JSON 格式的数据。
	w.Header().Set("Content-Type", "application/json")
	// 发送 HTTP 状态码 200 OK 到客户端，表示请求已成功处理。
	w.WriteHeader(http.StatusOK)
	h := H{
		Code: code,
		Data: data,
		Msg:  msg,
	}
	ret, err := json.Marshal(h)
	if err != nil {
		fmt.Println(err)
	}
	// 将序列化后的 JSON 数据 ret 写入响应流，发送给客户端。
	w.Write(ret)
}

func RespList(w http.ResponseWriter, code int, data interface{}, total interface{}) {
	// 设置响应的 Content-Type 头为 application/json，告诉客户端响应体是 JSON 格式的数据。
	w.Header().Set("Content-Type", "application/json")
	// 发送 HTTP 状态码 200 OK 到客户端，表示请求已成功处理。
	w.WriteHeader(http.StatusOK)
	h := H{
		Code:  code,
		Rows:  data,
		Total: total,
	}
	ret, err := json.Marshal(h)
	if err != nil {
		fmt.Println(err)
	}
	// 将序列化后的 JSON 数据 ret 写入响应流，发送给客户端。
	w.Write(ret)
}

func RespFail(w http.ResponseWriter, msg string) {
	Resp(w, -1, nil, msg)
}

func RespOK(w http.ResponseWriter, data interface{}, msg string) {
	Resp(w, 0, data, msg)
}

func RespOKList(w http.ResponseWriter, data interface{}, total interface{}) {
	RespList(w, 0, data, total)
}
