// this file is to interact with tbl_user_file table in database
// 用于操作用户文件表, 每一个用户都拥有自己的文件表
// 用户文件表存储了用户上传的文件的元信息
package db

import (
	"fileserver/fileserver/db/mysql"
	"fileserver/fileserver/orm"
	"fmt"
)

// TODO: insert or udpate 存储用户上传的文件元信息进入用户文件表
func Upload2UserFileDB(userfile *orm.UserFile) bool {
	// db prepare statment
	stmt, err := mysql.GetDBConnection().Prepare(
		"insert ignore into tbl_user_file(`user_name`, `file_sha1`, `file_size`, `file_name`) values(?,?,?,?)")
	if err != nil {
		fmt.Println("Failed to prepare statement, err: ", err.Error())
		return false
	}
	defer stmt.Close()

	// execute sql
	_, err = stmt.Exec(userfile.UserName.String,
		userfile.FileSha1.String,
		userfile.FileSize.Int64,
		userfile.FileName.String)

	if err != nil {
		fmt.Println("Failed to exec statement, err: ", err.Error())
		return false
	}

	return true
}

// TODO: select 从用户文件表中获取用户文件元信息
// QueryUserFileMetas: 根据username, limit限制返回的数据条数，返回用户文件元信息
func QueryUserFileMetas(username string, limit int) ([]orm.UserFile, error) {
	// db prepare statment
	return nil, nil
}

// TODO: delete 删除用户文件元信息
// 删除该用户的用户文件表元信息，但是不要删除用户上传的文件，该文件的元信息存储在tbl_file表中
func DeleteUserFile(username string, filehash string) bool {
	// db prepare statment
	return true
}
