package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"main/meta"
	"main/util"
	"net/http"
	"os"
	"time"
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

		//接收文件到本地目录，返回句柄和文件头
		//FormFile表示从上传请求中取出名叫 file 的那个上传字段
		file, head, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("Failed to get data,err: %v\n", err)
			return
		}

		//收尾关闭文件
		defer file.Close()

		//记录文件元信息
		filemeta := meta.FileMeta{
			FileName: head.Filename,
			Location: "./tmp/" + head.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		//创建新文件，返回句柄
		newFile, err := os.Create(filemeta.Location)
		if err != nil {
			fmt.Printf("Failed to create file,err: %v\n", err)
			return
		}

		//收尾关闭新文件句柄
		defer newFile.Close()

		//复制文件到指定路径
		filemeta.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("Failed to save data into file,err: %v\n", err)
			return
		}

		//将文件指针移动到文件头
		newFile.Seek(0, 0)
		filemeta.FileSha1 = util.FileSha1(newFile)
		meta.UpdateFileMeta(filemeta)

		//跳转到上传成功页面，code是302临时重定向
		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)

	}

}

// 上传完成路由函数
func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload finished.")
}

// 文件信息查询接口
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {

	//解析请求中的表单和URl参数
	r.ParseForm()

	//从请求参数中取得第一个filehash
	fileHash := r.Form["filehash"][0]
	fmeta := meta.GetFileMeta(fileHash)

	//把结构体转为json
	data, err := json.Marshal(fmeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("Failed to get filemeta,err: %v\n", err)
		return
	}

	w.Write(data)
}
