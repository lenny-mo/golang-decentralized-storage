// this file is used to store file meta data
package meta

import (
	"fileserver/fileserver/db"
	"sync"
)

// 文件元信息结构体
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string // 时间戳：上传时间经过格式化变成字符串
}

// 文件元信息结构体的map, 对外不可见
var fileMetas map[string]FileMeta // key: 文件的sha1值, value: 文件元信息

// 创建文件元信息map的互斥锁
var fileMetasMutex sync.Mutex

// init: 初始化map
func init() {
	// 初始化map
	fileMetas = make(map[string]FileMeta)
}

// UpdateFileMeta: 新增/更新文件元信息
func UpdateFileMeta(fmeta *FileMeta) {
	// 通过互斥锁保证并发安全
	fileMetasMutex.Lock()
	defer fileMetasMutex.Unlock()
	// 根据sha1值作为key, 文件元信息作为value, 更新map
	fileMetas[fmeta.FileSha1] = *fmeta
}

// UpdateFileMetaDB: add or update file meta to mysql, return true if success
func UpdateFileMetaDB(fmeta *FileMeta) bool {
	// thread safe
	return db.FileUploadFinished(fmeta.FileSha1,
		fmeta.FileName,
		fmeta.FileSize,
		fmeta.Location)
}

// GetFileMeta: 通过sha1值获取文件指针
func GetFileMeta(fileSha1 string) *FileMeta {

	data, ok := fileMetas[fileSha1]

	if !ok {
		// 如果没有找到对应的元信息，返回nil
		return nil
	}
	return &data
}

// GetFileMetaDB fetch filemeta from database using filehash
func GetFileMetaDB(filehash string) (*FileMeta, error) {
	tfile, err := db.GetFileMeta(filehash)
	if err != nil {
		return nil, err
	}
	// there could be a problem
	fmeta := FileMeta{
		FileSha1: tfile.FileHash,
		FileSize: tfile.FileSize.Int64,
		FileName: tfile.FileName.String,
		Location: tfile.FileAddr.String,
	}

	return &fmeta, nil
}

// RemoveFileMeta: 通过key删除map中的元素
func RemoveFileMeta(fileSha1 string) {
	// 通过互斥锁保证并发安全
	fileMetasMutex.Lock()
	defer fileMetasMutex.Unlock()
	delete(fileMetas, fileSha1)
}
