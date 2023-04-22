package db

import (
	"database/sql"
	"errors"
	"fileserver/fileserver/db/mysql"
	"fmt"
)

type Table struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}

// FileUploadFinished when file upload finished, insert file info into mysql
func FileUploadFinished(filehash string, filename string,
	filesize int64, fileaddr string) bool {

	// use prepare statement to avoid sql injection, and use ignore to avoid duplicate
	// stmt: a prepared statement which can be executed multiple times

	// get db connection
	dbConnect := mysql.GetDBConnection()
	// if the connection is nil, return false, this will not cause panic
	if dbConnect == nil {
		fmt.Println("Failed to connect to mysql")
		return false
	}

	stmt, err := dbConnect.Prepare(
		"insert ignore into tbl_file (`file_sha1`, `file_name`, `file_size`," +
			"`file_addr`, `status`) values (?, ?, ?, ?, 1)")

	if err != nil {
		fmt.Println("Failed to prepare statement, err: " + err.Error())
		return false
	}
	// close statment after use
	defer stmt.Close()

	ret, err := stmt.Exec(filehash, filename, filesize, fileaddr)

	if err != nil {
		fmt.Println("Failed to exec statement, err: " + err.Error())
		return false
	}

	// check if the file has been uploaded before, because if the file has been uploaded
	// the insert will be ignored and return 0
	// RowsAffected: return the number of rows affected by the statement,
	// if rf <= 0, the file has been uploaded before
	if rf, err := ret.RowsAffected(); err == nil {
		if rf <= 0 {
			fmt.Println("File with hash: " + filehash + " has been uploaded before")
		}
		return true
	}
	return false // 插入失败
}

// GetFileMeta get a file meta data from db
func GetFileMeta(filehash string) (*Table, error) {
	// connect to db and use prepare statement to fetch a file meta data
	dbConnect := mysql.GetDBConnection()
	if dbConnect == nil {
		fmt.Println("Failed to connect to mysql")
		return nil, errors.New("failed to connect to mysql")
	}

	// use prepare statement to fetch a filemeta, status = 1 means the file is valid
	stmt, err := dbConnect.Prepare("select file_sha1, file_name, file_size, file_addr " +
		"from tbl_file where file_sha1 = ? and status = 1 limit 1")
	if err != nil {
		fmt.Println("Failed to prepare fetch statement, err: " + err.Error())
	}

	// close statement after use
	defer stmt.Close()

	// create a table to store the file meta data
	t := Table{}
	// use QueryRow to fetch a row, if the file is not found, return nil
	err = stmt.QueryRow(filehash).Scan(&t.FileHash, &t.FileName, &t.FileSize, &t.FileAddr)
	if err != nil {
		fmt.Println("Failed to fetch file meta data, filehash: " + filehash)
		return nil, err
	}

	// if the file is not found, return nil
	return &t, nil
}
