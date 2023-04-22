package handler

import (
	"encoding/json"
	"fileserver/fileserver/meta"
	"fileserver/fileserver/util"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// UploadHandler: 处理文件上传
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "Internal server error")
			return
		}

		io.WriteString(w, string(data))

	} else if r.Method == "POST" {
		// file 是一个文件流, head 是文件头包含了文件名
		file, head, err := r.FormFile("file")

		if err != nil {
			fmt.Printf("Fail to get data, err: %s", err.Error())
			http.Error(w, "Fail to get data", http.StatusBadRequest)
			return
		}

		defer file.Close()

		// 创建filemeta 对象, 只需要创建一次，因为后面的filemeta都是对这个对象的引用
		filemeta := meta.FileMeta{
			FileName: head.Filename,
			Location: "./tmp/" + head.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"), // 创建时间并且格式化
		}

		// put filemeta into map, 此时map中的filemeta的filesize为0
		meta.UpdateFileMeta(&filemeta)

		// 创建本地文件接受上传的文件流
		newfile, err := os.Create(filemeta.Location)

		if err != nil {
			// (HTTP 500)，表明发生了内部服务器错误，服务端创建文件失败
			http.Error(w, "Failed to create file", http.StatusInternalServerError)
			return
		}

		defer newfile.Close()

		// 将file -> newfile; file is a source, newfile is a destination
		// filesize 是从file中读取的字节数
		filemeta.FileSize, err = io.Copy(newfile, file)
		if err != nil {
			fmt.Printf("Fail to save data into file, err:%s\n", err.Error())
			http.Error(w, "Fail to save data into file", http.StatusInternalServerError)
			return
		}

		// 将文件位置从文件开头移动到 0 字节，也就是将文件位置重置为文件的开头。
		newfile.Seek(0, 0)
		// sha1值是string类型
		filemeta.FileSha1 = util.FileSha1(newfile)
		fmt.Println("filemeta.FileSha1: ", filemeta.FileSha1)
		// 更新filemeta
		meta.UpdateFileMeta(&filemeta)
		// update filemeta into database
		_ = meta.UpdateFileMetaDB(&filemeta)

		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}
}

// UploadSuccessHandler 上传完成
func UploadSuccessHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload successed!")
}

// GetFileMetaHandler: 获取文件元信息
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm() // 解析参数，默认是不会解析的

	filehash := r.Form["filehash"][0] // filehash 是客户端传过来的参数 key, return []string

	// fmeta := meta.GetFileMeta(filehash)
	fmeta, err := meta.GetFileMetaDB(filehash)
	if err != nil {
		fmt.Println("GetFileMetaDB err: ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// using json Marshal to convert fmeta to json format
	data, err := json.Marshal(*fmeta) // 将fmeta转换成json格式返回给客户端
	if err != nil {
		// if err when covertion, return 500
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

// DownloadHandler：下载文件
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fsha1 := r.Form.Get("filehash")

	fm := meta.GetFileMeta(fsha1)

	f, err := os.Open(fm.Location) // 根据文件路径打开文件, 并且读取到内存中

	if err != nil {
		// 如果打开文件失败，返回500
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 如果文件打开失败，文件句柄 f 将为 nil，因此在这种情况下无需关闭文件
	defer f.Close()

	data, err := ioutil.ReadAll(f) // 读取文件内容到内存中，仅限于小文件，大文件需要使用文件流

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("content-disposition", "attachment; filename=\""+fm.FileName+"\"")
	w.Write(data)
}

// FileMetaUpdateHandler: 更新元文件名
func FileMetaUpdateHandler(w http.ResponseWriter, r *http.Request) {
	// 解析请求参数
	r.ParseForm()

	// 需要传递三个参数
	opType := r.Form.Get("op")            // op: 0 重命名
	fileSha1 := r.Form.Get("filehash")    // 获取文件hash
	newFileName := r.Form.Get("filename") // 获取新文件名

	if opType != "0" {
		// 如果不是重命名操作, 返回403，403表示禁止访问
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if r.Method != "POST" {
		// 如果不是post请求，返回405，405表示请求方法不被允许
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// 根据sha1值获取filemeta信息
	curFileMeta := meta.GetFileMeta(fileSha1)
	// 更新文件名
	curFileMeta.FileName = newFileName
	meta.UpdateFileMeta(curFileMeta)

	// 把更新后的filemeta转换成json格式返回给客户端
	data, err := json.Marshal(curFileMeta)
	if err != nil {
		// 如果转换失败，返回500，500表示内部服务器错误
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 返回200，200表示成功
	w.WriteHeader(http.StatusOK)
	// 返回json格式的数据
	w.Write(data)
}

func FileDeleteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// 获取文件hash
	fileSha1 := r.Form.Get("filehash")

	fMeta := meta.GetFileMeta(fileSha1)

	// 删除文件元信息
	meta.RemoveFileMeta(fileSha1)

	// 删除本地文件
	err := os.Remove(fMeta.Location)

	if err != nil {
		// 如果删除失败，返回500，500表示内部服务器错误
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 返回200，200表示成功
	w.WriteHeader(http.StatusOK)
}
