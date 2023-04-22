package db

import (
	"fileserver/fileserver/db/mysql"
	"fmt"
)

// UserSignUp sign up a user and insert into table
func UserSignUp(username string, password string) bool {
	// use prepare statement to avoid sql injection, and use ignore to avoid duplicate
	stmt, err := mysql.GetDBConnection().Prepare("insert ignore into tbl_user (`user_name`, `user_pwd`)" +
		" values (?, ?)")
	if err != nil {
		fmt.Println("Failed to prepare statement, err: " + err.Error())
		return false
	}

	defer stmt.Close()

	res, err := stmt.Exec(username, password)
	if err != nil {
		fmt.Println("Failed to exec statement, err: " + err.Error())
		return false
	}
	// check if the user has been signed up before
	if rf, err := res.RowsAffected(); err == nil {
		if rf <= 0 {
			fmt.Println("User: " + username + " has been signed up before")
			return false
		}
	}

	return true
}
