// Description: 工具类
// Date: 2017-12-28 15:00:00
package util

import (
	"crypto/sha1"
	"encoding/hex"
	"hash"
	"io"
	"os"
	"path/filepath"
)

// Sha1Stream: sha1流
type Sha1Stream struct {
	_sha1 hash.Hash
}

// Update: 更新sha1值
func (obj *Sha1Stream) Update(data []byte) {
	if obj._sha1 == nil {
		obj._sha1 = sha1.New()
	}
	obj._sha1.Write(data)
}

// GetFileSize: 获取文件大小
func GetFileSize(filename string) int64 {
	var result int64

	filepath.Walk(filename, func(path string, f os.FileInfo, err error) error {
		result = f.Size()
		return nil
	})

	return result
}

// FileSha1: 计算文件的sha1值
func FileSha1(file *os.File) string {
	_sha1 := sha1.New()                       // 创建sha1对象
	io.Copy(_sha1, file)                      // 将file -> _sha1
	return hex.EncodeToString(_sha1.Sum(nil)) // 返回sha1值
}

// Sha1 create a sha1 hash with given string
func Sha1(str string) string {
	_sha1 := sha1.New()
	_sha1.Write([]byte(str))
	return hex.EncodeToString(_sha1.Sum(nil))
}
