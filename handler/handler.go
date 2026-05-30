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

		fmt.Printf("Sha1 of the %s is %s\n", filemeta.FileName, filemeta.FileSha1)

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
	fileHash := r.Form.Get("filehash")
	fmeta := meta.GetFileMeta(fileHash)

	//把结构体转为json
	data, err := json.Marshal(fmeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("Failed to get filemeta,err: %v\n", err)
		return
	}

	//返回响应
	w.Write(data)
}

// 文件下载接口
func DownloadHandler(w http.ResponseWriter, r *http.Request) {

	//解析参数
	r.ParseForm()
	fsha1 := r.Form.Get("filehash")

	//获取文件位置并打开
	fmeta := meta.GetFileMeta(fsha1)
	file, err := os.Open(fmeta.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("Failed to get target file,err: %v\n", err)
		return
	}

	//收尾关闭句柄
	defer file.Close()

	//读取全部文件内容
	data, err := ioutil.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("Failed to read target file,err: %v\n", err)
		return
	}

	//将文件内容作为下载响应返回给服务器
	//二进制数据流，不作为页面显示
	w.Header().Set("Content-Type", "application/octet-stream")
	//附件触发下载
	w.Header().Set("Content-Disposition", "attachment;filename=\""+fmeta.FileName+"\"")
	//返回响应
	w.Write(data)
}

// 文件重命名接口
func FileMetaUpdateHandler(w http.ResponseWriter, r *http.Request) {
	//解析参数
	r.ParseForm()

	//获取参数
	opType := r.Form.Get("op")
	filesha1 := r.Form.Get("filehash")
	newFileName := r.Form.Get("filename")

	//操作码不符
	if opType != "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	//操作不符
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	//重命名
	curFileMeta := meta.GetFileMeta(filesha1)
	curFileMeta.FileName = newFileName
	meta.UpdateFileMeta(curFileMeta)

	//返回json响应
	data, err := json.Marshal(curFileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("Failed to get filemeta,err: %v\n", err)
		return
	}
	w.Write(data)

}

// 删除文件及元信息接口
func FileDeleteHandler(w http.ResponseWriter, r *http.Request) {
	//解析参数
	r.ParseForm()

	//拿到参数
	filesha1 := r.Form.Get("filehash")

	//删除磁盘文件
	fmeta := meta.GetFileMeta(filesha1)
	os.Remove(fmeta.Location)

	//移除元信息并响应
	meta.RemoveFileMeta(filesha1)
	w.WriteHeader(http.StatusOK)
}
