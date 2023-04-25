// this file is used to store user info and interacts with the user info in database
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

// UserSignin sign in a user with name and password
// return true if the user exists and the password is correct
func UserSignin(username string, encpwd string) bool {
	// using prepare statement to run sql query
	stmt, err := mysql.GetDBConnection().Prepare("select * from tbl_user where user_name=? limit 1")
	if err != nil {
		fmt.Println("Failed to prepare statement, err: " + err.Error())
		return false
	}

	defer stmt.Close() // close the statement after use

	rows, err := stmt.Query(username)
	if err != nil {
		fmt.Println("Failed to query statement, err: " + err.Error())
		return false
	}
	defer rows.Close() // close the rows since they save the result and hold the connection

	//check if the user exists
	if !rows.Next() {
		fmt.Println("Username:", username, "does not exist")
	}

	// parse the rows
	pRows, _ := mysql.ParseUserRows(rows)

	if len(pRows) > 0 && pRows[0].UserPwd.Valid {
		return pRows[0].UserPwd.String == encpwd
	}

	return false
}

// UpdateToken update user token into db
func UpdateToken(username, token string) bool {
	// connect to db and using db to prepare statement
	// the prepare statement can be executed multiple times
	stmt, err := mysql.GetDBConnection().Prepare("replace into tbl_user_token (`user_name`, `user_token`)" +
		" values (?, ?)")

	if err != nil {
		fmt.Println("Failed to prepare statement, err: " + err.Error())
		return false
	}

	// since the statement is a resource, it should be closed after use
	defer stmt.Close()

	// execute the statement
	_, err = stmt.Exec(username, token)

	if err != nil {
		fmt.Println("Failed to exec statement, err: " + err.Error())
		return false
	}

	// if the statement is executed successfully, return trued
	return true
}
