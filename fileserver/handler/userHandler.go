// this file is used to store the user information in the database
package handler

import (
	"crypto/rand"
	"encoding/base64"
	"fileserver/fileserver/db"
	"fileserver/fileserver/util"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	// password_salt is a random string used to be added in the password
	password_salt = generateSalt(16)
	// token_salt is a random string used to be added in the token
	token_salt = generateSalt(16)
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
	// TODO: store email into database
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	// check if the username and password are correct
	if len(username) < 3 || len(password) < 6 {
		w.Write([]byte("Invalid parameter"))
		return
	}

	// encrypt the password with salt
	encodePassword := util.Sha1(password + password_salt)
	// insert user info into mysql
	ok := db.UserSignUp(username, encodePassword)
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
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/static/view/signin.html", http.StatusFound)
		return
	}

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

	// if the user exists, generate a token, the token can be used for 1 hour

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
// the length of the token is 64
// TODO: the length of token should be 64
func genearteToken(username string) string {
	timeStamp := fmt.Sprintf("%x", time.Now().Unix())
	fmt.Println("timestamp of the user", timeStamp)

	tokenPrefix := util.MD5(username + timeStamp + token_salt)

	return tokenPrefix + "_" + timeStamp[:8]
}
