// 这个文件处理用户的登陆状态
package session

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/boj/redistore"
)

var (
	store *redistore.RediStore
)

func init() {
	//最大空闲连接数10，用于加密和解密存储在Redis中的会话数据的秘钥·······
	var err error
	fmt.Println("start to connect to redis")
	store, err = redistore.NewRediStore(10, "tcp", "localhost:6379", "", []byte("redis-secret-key"))

	if err != nil {
		fmt.Println("cannot connect to redis")
	} else {
		fmt.Println("successfully connect to redis")
	}
}

type UserSession struct {
	Username string
	Token    string
}

// GetSessionUser 获取用户的session信息，如果没有在redis中找到用户的session信息，返回nil
func GetSessionUser(r *http.Request) *UserSession {
	session, _ := store.Get(r, "user")
	s, ok := session.Values["user"]
	if !ok {
		return nil
	}
	u := &UserSession{}
	json.Unmarshal([]byte(s.(string)), u)
	return u
}

// SaveSessionUser 保存用户的session信息到redis中
func SaveSessionUser(w http.ResponseWriter, r *http.Request, u *UserSession) {
	session, _ := store.Get(r, "user")
	data, _ := json.Marshal(u)
	session.Values["user"] = string(data)
	store.Save(r, w, session)
}

// ClearSessionUser 清除用户的session信息
func ClearSessionUser(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "user")
	session.Options.MaxAge = -1 // 将会话标记为已过期，这将导致会话在下一次请求期间被删除
	store.Save(r, w, session)   // 将更新的会话保存回会话存储。由于 MaxAge 已设置为 -1，因此当响应发送到客户端时，会话将被删除。
}
