package handler

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

// 上传路由函数，无需手动传参，w为响应写入器，r为客户端请求
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		//返回上传页面
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "Internel server error")
			return
		}
		//没错误直接返回index.html
		//WriteString用来返回内容
		io.WriteString(w, string(data))
	} else if r.Method == "POST" {
		//接收文件到本地目录
		//FormFile表示从上传请求中取出名叫file的文件
		file, head, err := r.FormFile("file")

		if err != nil {
			fmt.Printf("Failed to get data,err: %v\n", err)
			return
		}

		//收尾关闭文件
		defer file.Close()

		newFile, err := os.Create("./tmp/" + head.Filename)
		if err != nil {
			fmt.Printf("Failed to create file,err: %v\n", err)
			return
		}

		defer newFile.Close()

		_, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("Failed to save data into file,err: %v\n", err)
			return
		}

		//跳转到上传成功页面，code是302临时重定向
		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)

	}

}

// 上传完成路由函数
func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload finished.")
}
