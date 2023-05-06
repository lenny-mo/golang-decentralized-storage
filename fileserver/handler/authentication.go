// this file is to check if the user token is valid
// 提供身份验证相关的操作
package handler

import (
	"fileserver/fileserver/session"
	"net/http"
)

// HTTPInterceptor is a middleware to check if the user token is valid
// 用于验证用户是否登陆过，如果没有登陆过，跳转到登陆页面
func SessionAuthInterceptor(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// Check the user session
			userSession := session.GetSessionUser(r)

			// If userSession is nil, the user is not authenticated
			if userSession == nil {
				// Redirect to the login page
				http.Redirect(w, r, "/user/signin", http.StatusFound)
				return
			}

			// The user is authenticated, call the original handler
			h(w, r)
		},
	)
}
