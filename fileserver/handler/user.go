package handler

import (
	"crypto/rand"
	"encoding/base64"
	"fileserver/fileserver/db"
	"fileserver/fileserver/util"
	"io/ioutil"
	"net/http"
)

var (
	password_salt = generateSalt(16)
)

// SignUpHandler handle user sign up request
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
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
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	// check if the username and password are correct
	if len(username) < 3 || len(password) < 6 {
		w.Write([]byte("Invalid parameter"))
		return
	}

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
