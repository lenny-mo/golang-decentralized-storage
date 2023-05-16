package mysql

import (
	"database/sql"
	"fileserver/fileserver/orm"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB // global sql connection

func init() {
	// define err to store error message
	var err error

	// cannot use short variable declaration here, because db is a global variable
	db, err = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/fileserver?charset=utf8")

	if err != nil {
		// if connection failed, panic and print error message to console
		fmt.Println("Failed to connect to mysql, err: " + err.Error())
	}

	// set max open connections to 1000, default is 0 (unlimited)
	db.SetMaxOpenConns(1000)

	// try to use ping method to check if the connection is ok
	err = db.Ping()

	if err != nil {
		// if connection failed, panic and print error message to console
		fmt.Println("Failed to connect to mysql, err: " + err.Error())
		// force the program to exit
		os.Exit(1)
	} else {
		fmt.Println("Successfully ping to master mysql")
	}
}

// GetDBConnection: return the global sql connection
func GetDBConnection() *sql.DB {
	// if db is nil, then print error message and exit
	return db
}

// ParseRows parse the rows to file meta data
func ParseUserRows(rows *sql.Rows) ([]*orm.UserInfo, error) {
	// create a slice to store the file meta data
	userinfo := make([]*orm.UserInfo, 0)

	// iterate the rows
	for rows.Next() {
		// create a user info struct
		u := &orm.UserInfo{}

		// scan the row and store the data to file meta data
		err := rows.Scan(&u.UserName, &u.UserPwd)
		if err != nil {
			fmt.Println("Failed to scan row, err: " + err.Error())
			return nil, err
		}

		// append the file meta data to the slice
		userinfo = append(userinfo, u)
	}

	return userinfo, nil
}
