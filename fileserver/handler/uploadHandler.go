// 这个文件的作用是处理文件上传的请求
package handler

import (
	"database/sql"
	"encoding/json"
	"fileserver/fileserver/db"
	"fileserver/fileserver/meta"
	"fileserver/fileserver/orm"
	"fileserver/fileserver/session"
	"fileserver/fileserver/util"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// UploadHandler: 处理文件上传, 如果文件大于100MB, 触发分块上传
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		data, err := os.ReadFile("./static/view/index.html") // 读取整个文件内容
		if err != nil {
			io.WriteString(w, "Internal server error")
			return
		}

		io.WriteString(w, string(data)) // 把请求的页面内容写入到w中

	} else if r.Method == "POST" {
		// file 该变量保存上传文件的文件流
		// head 保存了关于文件的元数据，比如文件名、内容类型等。它的类型是*multipart.FileHeader
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
			UploadAt: time.Now().Format("2006-01-02 15:04:05"), // 使用称为“Unix 日期”的参考时间来表示所需日期时间字符串的布局
		}

		// put filemeta into map, 此时map中的filemeta的filesize为0
		meta.UpdateFileMeta(&filemeta)

		// 创建本地文件 接受上传的文件流
		newfile, err := os.Create(filemeta.Location)
		if err != nil {
			// (HTTP 500)，表明发生了内部服务器错误，服务端创建文件失败
			http.Error(w, "Failed to create file", http.StatusInternalServerError)
			return
		}
		defer newfile.Close()

		// filesize 复制的字节数，已传输文件的大小
		filemeta.FileSize, err = io.Copy(newfile, file)
		if err != nil {
			fmt.Printf("Fail to save data into file, err:%s\n", err.Error())
			http.Error(w, "Fail to save data into file", http.StatusInternalServerError)
			return
		}

		// 将文件指针从文件开头移动到 0 字节，也就是将文件位置重置为文件的开头。
		newfile.Seek(0, 0)
		// 创建一个用于计算文件sha1值的channel，缓冲为零
		sha1Chan := make(chan string)
		go func() {
			// 计算文件的sha1值
			sha1Chan <- util.FileSha1(newfile)
		}()

		// 从channel中读取sha1值
		filemeta.FileSha1 = <-sha1Chan
		fmt.Println("filemeta.FileSha1: ", filemeta.FileSha1)
		// update filemeta
		meta.UpdateFileMeta(&filemeta)
		// update filemeta into tbl_file database
		go func() {
			_ = meta.UpdateFileMetaDB(&filemeta)
		}()

		// get the user and file info to update tbl_user_file
		s := session.GetSessionUser(r)
		username := s.Username

		u := &orm.UserFile{UserName: sql.NullString{String: username, Valid: true}}
		u.FileSha1 = sql.NullString{String: filemeta.FileSha1, Valid: true}
		u.FileName = sql.NullString{String: filemeta.FileName, Valid: true}
		u.FileSize = sql.NullInt64{Int64: filemeta.FileSize, Valid: true}

		// 上传文件的同时，更新该用户的 tbl_user_file，用户每次上传文件都会更新该表，不管文件是否重复
		go func() {
			_ = db.Upload2UserFileDB(u)
		}()

		// 重定向到用户的home页面
		http.Redirect(w, r, "/user/info", http.StatusFound)
	}
}

// UploadSuccessHandler this func is deprecated
// if the user upload file successfully, redirect to home page rather than this
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

// FileMetaUpdateHandler:
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

// TODO: 查询批量的用户文件表信息 from tbl_user_file
// QueryUserFileMetas : 查询批量的用户文件表信息 from tbl_user_file
func QueryUserFileMetas(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	limit := r.Form.Get("limit")

	fmt.Println("username: ", username, "limit: ", limit)
	// 根据username和limit查询用户文件表信息, 并且返回给客户端
}
