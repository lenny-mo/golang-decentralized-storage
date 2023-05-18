// 管理与redis的连接池
package redis

import "github.com/go-redis/redis"

var (
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6380", // connect to localhost
		Password: "",               // default no password
		DB:       0,                // use default DB
		PoolSize: 100,              // 连接池大小, 默认是10,
	})
)	

func NewRedisClient() *redis.Client {  
	return redisClient
}

func CloseRedisClient() {
	redisClient.Close()
}

// 测试redis连接是否正常
func Ping() error {
	return redisClient.Ping().Err()
}
