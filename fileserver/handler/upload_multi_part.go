// 这个文件处理分块上传
package handler

import (
	"fileserver/fileserver/cache/redis"
	"fmt"
	"math"
	"net/http"
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
	// 2. 读取用户上传的分块

	r.ParseForm()

	// 获取redis链接
	// 获取文件handle, 用于存储分块内容
	// 更新redis缓存中的分块信息
	// 返回处理结果给客户端
}

// TODO: 上传完成
func CompleteUploadMultiPartHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 解析用户请求
}

// TODO: 通知上传合并
func UploadCombineHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 解析用户请求
	// 获取redis链接
	// 判断是否所有分块都上传完成
	// 合并分块
	// 更新 tbl_file 表
	// 更新 tbl_user_file 表
	// 删除redis缓存中的分块信息
	// 返回处理结果给客户端
}

// TODO: 取消上传
func CancelUploadHandler(w http.ResponseWriter, r *http.Request) {
	// 删除已经存在的分块
	// 删除redis缓存中的分块信息
	// 获取redis链接
	// 删除redis缓存中的分块信息
	// 返回处理结果给客户端
}

// TODO: 查询分块上传的状态
// 对应的router /file/mpupload/status
func QueryUploadStatusHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 解析用户请求
	// 2. 读取用户上传的分块
	// 获取redis链接
	// 获取redis缓存中的分块信息
	// 返回处理结果给客户端
}
