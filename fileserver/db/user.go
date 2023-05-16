// this file is used to store and get user info, and interacts with the user info in database
package db

import (
	"fileserver/fileserver/db/mysql"
	"fileserver/fileserver/orm"
	"fmt"
)

// UserSignUp sign up a user and insert into table tbl_user
func UserSignUp(username, password, email string) bool {
	// use prepare statement to avoid sql injection, and use ignore to avoid duplicate
	// only prevents insertion of rows that would cause a duplicate key value in a unique index or primary key
	// 已经删除user_name的unique key 约束
	stmt, err := mysql.GetDBConnection().Prepare("insert into tbl_user (`user_name`, `user_pwd`, `email`)" +
		" values (?, ?, ?)")
	if err != nil {
		fmt.Println("Failed to prepare statement, err: " + err.Error())
		return false
	}

	defer stmt.Close()

	res, err := stmt.Exec(username, password, email)
	if err != nil {
		fmt.Println("Failed to exec statement, err: " + err.Error())
		return false
	}

	// check if the user has been signed up before
	// if the user signed up before, the rows affected should be 0
	if rf, err := res.RowsAffected(); err == nil {
		if rf <= 0 {
			fmt.Println("rows affected is: ", rf)
			fmt.Println("User: " + username + " has been signed up before")
			return false
		}
		fmt.Println("rows affected is: ", rf)
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
	// check if the query is successful
	if err != nil {
		fmt.Println("Failed to query statement, err: " + err.Error())
		return false
	}
	defer rows.Close() // close the rows since they save the result and hold the connection

	//check if the user exists
	if !rows.Next() {
		fmt.Println("Username:", username, "does not exist")
		return false
	}

	// parse the rows
	pRows, _ := mysql.ParseUserRows(rows)

	if len(pRows) > 0 && pRows[0].UserPwd.Valid {
		fmt.Println("User password is: ", pRows[0].UserPwd.String)
		return pRows[0].UserPwd.String == encpwd
	}

	return true
}

// UpdateToken update user token into db
func UpdateToken(username, token string) bool {
	// connect to db and using db to prepare statement
	// the prepare statement can be executed multiple times
	// use replace since we want to update the token if the user has signed in before
	// if there is a token before, it will be replaced with the new token
	// we must update token with an expired time, such as 1 hour
	stmt, err := mysql.GetDBConnection().Prepare("replace into user_token (`user_name`, `token`, `expires_at`)" +
		" values (?, ?, DATE_ADD(NOW(), INTERVAL 1 HOUR))")

	// if replace failed, return false
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

// GetUserInfo get user info from db
func GetUserInfo(username string) (*orm.UserInfo, error) {
	// using prepare statement to run sql query
	stmt, err := mysql.GetDBConnection().Prepare("select user_name," +
		" email," +
		" phone," +
		" signup_at," +
		" last_active," +
		" profile," +
		" status" +
		" from tbl_user where user_name=? limit 1")

	if err != nil {
		fmt.Println("Failed to prepare statement, err: " + err.Error())
		return nil, err
	}

	defer stmt.Close() // close the statement after use

	//
	u := orm.UserInfo{}
	err = stmt.QueryRow(username).Scan(&u.UserName, &u.Email, &u.Phone, &u.SignupAt, &u.LastActive, &u.Profile, &u.Status)

	if err != nil {
		fmt.Println("Failed to query statement, err: " + err.Error())
		return &u, err
	}

	return &u, nil
}
