package main

import (
	"fmt"
	"main/handler"
	"net/http"
)

func main() {
	//注册路由，将上传转交给Upload_handler处理
	http.HandleFunc("/file/upload", handler.Upload_handler)
	//启动服务器，监听8080
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Fail to start the server,err:%v\n", err)
	}
}
