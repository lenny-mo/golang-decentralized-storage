// 这个文件处理分块上传和断点续传的逻辑
package handler

import "net/http"

// TODO: 实现初始化分块上传，并且返回分块上传的信息
func InituploadMultiPartHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 解析用户请求
	// 2. 读取用户上传的分块
	// 获取redis链接
	// 生成分块上传的信息
	// 将分块信息写入到redis缓存中
	// 4. 返回分块上传的信息给客户端

	r.ParseForm()
}

// TODO: 执行分块上传
func UploadMultiPartHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 解析用户请求
	// 2. 读取用户上传的分块

	r.ParseForm()

	// 获取redis链接
	// 获取文件handle, 用于存储分块内容
	// 更新redis缓存中的分块信息
	// 返回处理结果给客户端
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
func QueryUploadHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 解析用户请求
	// 2. 读取用户上传的分块
	// 获取redis链接
	// 获取redis缓存中的分块信息
	// 返回处理结果给客户端
}
