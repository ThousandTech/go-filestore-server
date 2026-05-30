package handler

import (
	"io"
	"io/ioutil"
	"net/http"
)

// 上传
func Upload_handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		//返回上传页面
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "Internel server error")
			return
		}
		io.WriteString(w, string(data))
	} else if r.Method == "POST" {
		//接收文件到本地
	}

}
