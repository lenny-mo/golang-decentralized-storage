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
	UserName sql.NullString
	UserPwd  sql.NullString
	Email    sql.NullString
}
