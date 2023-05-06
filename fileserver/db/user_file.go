// this file is to interact with tbl_user_file table in database
// 用于操作用户文件表, 每一个用户都拥有自己的文件表
// 用户文件表存储了用户上传的文件的元信息
package db

import (
	"database/sql"
)

// orm object
type UserFile struct {
	ID         sql.NullInt32
	UserName   sql.NullString
	FileSha1   sql.NullString
	FileSize   sql.NullInt64
	FileName   sql.NullString
	UploadAt   sql.NullTime
	LastUpdate sql.NullTime
	Status     sql.NullInt32
}

// TODO: insert or udpate 存储用户上传的文件元信息进入用户文件表
func Upload2UserFileDB(username string) bool {
	// db prepare statment

	return true
}

// TODO: select 从用户文件表中获取用户文件元信息
// QueryUserFileMetas: 根据username, limit限制返回的数据条数，返回用户文件元信息
func QueryUserFileMetas(username string, limit int) ([]UserFile, error) {
	// db prepare statment
	return nil, nil
}

// TODO: delete 删除用户文件元信息
// 删除该用户的用户文件表元信息，但是不要删除用户上传的文件，该文件的元信息存储在tbl_file表中
func DeleteUserFile(username string, filehash string) bool {
	// db prepare statment
	return true
}
