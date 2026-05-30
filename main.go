package main

import (
	"fmt"
	"main/handler"
	"net/http"
)

func main() {
	//注册路由，将上传转交给Upload_handler处理
	http.HandleFunc("/file/upload", handler.UploadHandler)
	//上传成功路由
	http.HandleFunc("/file/upload/suc", handler.UploadSucHandler)
	//获取文件信息路由
	http.HandleFunc("/file/meta", handler.GetFileMetaHandler)
	//文件下载路由
	http.HandleFunc("/file/download", handler.DownloadHandler)
	//文件重命名路由
	http.HandleFunc("/file/update", handler.FileMetaUpdateHandler)
	//文件删除路由
	http.HandleFunc("/file/delete", handler.FileDeleteHandler)

	//启动服务器，监听8080
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Fail to start the server,err: %v\n", err)
	}
}
