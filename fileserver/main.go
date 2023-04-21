package main

// import some packages
import (
	"fileserver/fileserver/handler"
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("main")
	// 使用http.HandleFunc()函数注册路由
	http.HandleFunc("/file/upload", handler.UploadHandler)

	http.HandleFunc("/file/upload/suc", handler.UploadSuccessHandler)

	http.HandleFunc("/file/meta", handler.GetFileMetaHandler)

	http.HandleFunc("/file/download", handler.DownloadHandler)

	http.HandleFunc("/file/update", handler.FileMetaUpdateHandler)

	http.HandleFunc("/file/delete", handler.FileDeleteHandler)

	fmt.Println("Start server at 8080")
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Printf("Failed to start server, err: %s", err.Error())
	}

}
