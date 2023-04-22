package mysql

import (
	"database/sql"
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
		fmt.Println("Successfully ping to mysql")
	}
}

// GetDBConnection: return the global sql connection
func GetDBConnection() *sql.DB {
	// if db is nil, then print error message and exit
	return db
}
