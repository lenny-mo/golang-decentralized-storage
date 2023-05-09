// 这个文件处理分块上传
package handler

import (
	"fileserver/fileserver/cache/redis"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"
)

// MultipartUploadInfo : 分块上传的信息
type MultipartUploadInfo struct {
	FileHash   string // 文件hash
	FileSize   int    // 文件大小
	UploadID   string // 上传的id, 即使文件重复上传，也会有不同的id
	ChunkSize  int    // 分块大小，注意最后一
	ChunkCount int    // 分块数量
}

// 实现初始化分块上传，并且返回分块上传的信息
// 对应的router /file/mpupload/init
func InitUploadMultiPartHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 解析用户请求
	r.ParseForm()
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filesize, err := strconv.Atoi(r.Form.Get("filesize")) // convert
	if err != nil {
		w.Write([]byte("params invalid"))
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// 2 connect to redis
	redisClient := redis.NewRedisClient()
	defer redisClient.Close()

	// 3 初始化分块信息
	chunksize := 5 * 1024 * 1024
	uploadInfo := &MultipartUploadInfo{
		FileHash:   filehash,
		FileSize:   filesize,
		UploadID:   username + fmt.Sprintf("%x", time.Now().UnixNano()),
		ChunkSize:  chunksize,                                                // 5MB
		ChunkCount: int(math.Ceil(float64(filesize) / (float64(chunksize)))), // Ceil 向上取整
	}

	// 4 将初始化信息写入到redis缓存中
	redisClient.MSet("MP_"+uploadInfo.UploadID+"_chunkcount", uploadInfo.ChunkCount,
		"MP_"+uploadInfo.UploadID+"_filehash", uploadInfo.FileHash,
		"MP_"+uploadInfo.UploadID+"_filesize", uploadInfo.FileSize)

	// 5 返回初始化信息给客户端
	w.Write([]byte("OK"))
}

// TODO: 执行分块上传
// 对应的router /file/mpupload/uppart
func UploadMultiPartHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 解析用户请求

	r.ParseForm()
	uploadID := r.Form.Get("uploadid")
	chunkIndex := r.Form.Get("index")

	// 获取redis链接
	redisClient := redis.NewRedisClient()
	defer redis.CloseRedisClient()

	fileHandler, err := os.Create("./tmp/" + uploadID + "/" + chunkIndex)
	if err != nil {
		w.Write([]byte("Upload part failed."))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Upload part failed: ", err)
		return
	}
	defer fileHandler.Close()

	// 将分块文件写入到本地, 并且从请求体中读取分块内容
	buffersize := 1024 * 1024
	buffer := make([]byte, buffersize)
	for {
		// TODO: 写入分块到文件的时候，需要判断这个分块的hash是否和用户上传的hash一致
		n, err := r.Body.Read(buffer)
		fileHandler.Write(buffer[:n])
		if err != nil {
			break
		}
	}

	// 更新redis缓存中的分块信息
	redisClient.SAdd("MP_"+uploadID+"_chunks", chunkIndex)

}

// TODO: 上传完成
func CompleteUploadMultiPartHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 解析用户请求
}

// TODO: 通知上传合并
func UploadCombineHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 解析用户请求
	r.ParseForm()
	uploadID := r.Form.Get("uploadid")

	redisClient := redis.NewRedisClient()

	// 2. 通过uploadid查询redis缓存，判断分块是否全部上传完成
	data, err := redisClient.MGet("MP_"+uploadID+"_chunkcount", "MP_"+uploadID+"_chunks").Result()
	if err != nil {
		w.Write([]byte("Upload combine failed"))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Upload combine failed: ", err)
		return
	}

	// 3. 判断分块是否全部上传完成
	totalCount, _ := strconv.Atoi(data[0].(string)) // 分块总数
	completeCount := len(data[1].([]interface{}))   // 已经上传的分块数量
	if totalCount != completeCount {
		w.Write([]byte("Upload combine failed"))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Upload combine failed: ", err)
		return
	}

	// TODO: 合并分块

	// 更新tbl_file and tbl_user_file
	// 如果文件已经存在，只需要更新 tbl_user_file
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filesize := r.Form.Get("filesize")
	filename := r.Form.Get("filename")

}

// TODO: 取消上传
func CancelUploadHandler(w http.ResponseWriter, r *http.Request) {

}

// TODO: 查询分块上传的状态
// 对应的router /file/mpupload/status
func QueryUploadStatusHandler(w http.ResponseWriter, r *http.Request) {

}
