package main

// import some packages
// for example, fmt
import (
	"fileserver/fileserver/handler"
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("main")

	http.HandleFunc("/file/upload", handler.SessionAuthInterceptor(handler.UploadHandler))

	http.HandleFunc("/file/upload/suc", handler.UploadSuccessHandler)

	http.HandleFunc("/file/meta", handler.GetFileMetaHandler)

	http.HandleFunc("/file/download", handler.DownloadHandler)

	http.HandleFunc("/file/update", handler.FileMetaUpdateHandler)

	http.HandleFunc("/file/delete", handler.FileDeleteHandler)

	http.HandleFunc("/user/signup", handler.SignUpHandler)

	http.HandleFunc("/user/signin", handler.SignInHandler)

	// user home page
	http.HandleFunc("/user/info", handler.UserInfoHandler)

	// TODO: 检验快速上传的router
	http.HandleFunc("/file/fastupload", handler.FastUploadHandler)

	// TODO: 分块上传的router

	// static resource
	// fs := http.FileServer(http.Dir("./static"))
	// http.Handle("/static/", http.StripPrefix("/static", fs))

	fmt.Println("Start server at 8080")

	// 启动server
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Printf("Failed to start server, err: %s", err.Error())
	}

}
