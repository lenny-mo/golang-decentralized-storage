// this file is used to calculate something about crypto
package util

import (
	"crypto/md5"
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

// FileSha1 create a sha1 hash with given file
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

// MD5 create a md5 hash with length 32
// you should know that MD5  is not secure enough compared with SHA1
// MD5 generate 128bit hash, SHA1 generate 160bit hash
func MD5(str string) string {
	_md5 := md5.New()
	_md5.Write([]byte(str))
	// convert into hex string
	return hex.EncodeToString(_md5.Sum(nil))
}
