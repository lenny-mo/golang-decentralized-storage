// this file is used to store the user information in the database
package handler

import (
	"crypto/rand"
	"encoding/base64"
	"fileserver/fileserver/db"
	"fileserver/fileserver/session"
	"fileserver/fileserver/util"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	// password_salt is a random string used to be added in the password
	// password_salt = generateSalt(16)
	password_salt = "1234567890123456"
	// token_salt is a random string used to be added in the token
	// token_salt = generateSalt(16)
	token_salt = "1234567890123456"
)

// SignUpHandler handle user sign up request
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	// if the request method is GET, then return the sign up page
	if r.Method == http.MethodGet {
		data, err := ioutil.ReadFile("./static/view/signup.html")
		if err != nil {
			// when read file failed, return 500
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal server error" + err.Error()))
			return
		}
		w.Write(data)
		return
	}

	// if the method is not GET, then it must be POST
	// parse the userdata from form
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	email := r.Form.Get("email")

	// check if the username and password are correct
	if len(username) < 3 || len(password) < 6 {
		w.Write([]byte("Invalid parameter, your username should be at least 3 characters and your password should be at least 6 characters"))
		return
	}

	// encrypt the password with salt
	encodePassword := util.Sha1(password + password_salt)
	// insert user info into mysql
	ok := db.UserSignUp(username, encodePassword, email)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to write into database"))
		return
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Sign up success!!!"))
	}

}

// SignInHandler handle user sign in request
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	// if the method is GET, then return the sign in page
	if r.Method == http.MethodGet {
		data, err := ioutil.ReadFile("./static/view/signin.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal server error" + err.Error()))
		}
		w.Write(data)
		return
	}

	// if the method is POST, then parse the form data
	// parse the form data, get the username and password
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	// check if the username and password are correct
	encryptedPassword := util.Sha1(password + password_salt)
	// check if the user exists
	ok := db.UserSignin(username, encryptedPassword)
	// if the user not exists, return 403, client forbidden
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Failed to sign in"))
		return
	}

	// generate a token for the user
	// if the user exists, generate a token, the token can be used for 1 hour
	token := genearteToken(username)
	// store the token into redis
	user := session.UserSession{Username: username, Token: token}
	session.SaveSessionUser(w, r, &user)

	// 重定向到file/upload
	http.Redirect(w, r, "/file/upload", http.StatusFound)
}

// TODO: 重定向到home.html, 并且能够显示用户上传的所有文件，也就是查询 tbl_user_file 表
// UserInfoHandler get the user info from the database, redirect to the home page
func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	// 获取用户的session info
	username := session.GetSessionUser(r).Username
	token := session.GetSessionUser(r).Token

	data := struct {
		Username string
		Token    string
	}{
		Username: username,
		Token:    token,
	}

	//TODO: 从数据库中查询用户的 tbl_user_file 表, 并且把用户上传过的所有文件返回给客户端

	t, err := template.ParseFiles("./static/view/home.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error" + err.Error()))
		return
	}

	// Render the template with the user data
	err = t.Execute(w, data)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error" + err.Error()))
		return
	}
}

// genearteSalt generate a random salt
func generateSalt(length int) string {
	// generate a random salt
	salt := make([]byte, length)
	// using random bytes to fill the salt
	_, err := rand.Read(salt)
	if err != nil {
		return "fileserver"
	}

	// encode the salt to base64 string
	return base64.StdEncoding.EncodeToString(salt)
}

// generateToken generate a token for the user,
// the length of the token is 40
func genearteToken(username string) string {
	timeStamp := fmt.Sprintf("%x", time.Now().Unix())
	fmt.Println("timestamp of the user", timeStamp)

	// the length of tokenPrefix is 32
	tokenPrefix := util.MD5(username + timeStamp + token_salt)

	return tokenPrefix + timeStamp[:8]
}

// TODO: complete this function
func IsTokenValid(token, username string) bool {
	// check if the token is expired in the database
	// check the time stamp of the token,
	// if the token is expired return false
	return len(token) == 40

	// fetch the token from the database and compare these two tokens are the same
}
