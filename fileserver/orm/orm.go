package orm

import (
	"database/sql"
)

// FileInfo is the struct of the table tbl_file
type FileInfo struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}

// UserInfo is the struct of the table tbl_user
type UserInfo struct {
	UserName   sql.NullString
	UserPwd    sql.NullString
	Email      sql.NullString
	Phone      sql.NullString
	SignupAt   sql.NullTime
	LastActive sql.NullTime
	Profile    sql.NullString
	Status     sql.NullInt32
}

// UserFile is the struct of the table tbl_user_file
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
