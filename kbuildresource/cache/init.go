package cache

import (
	"github.com/go-redis/redis"
	"time"
)

var (
	RedisClient *redis.Client
)

func init() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		Password: "",
		DialTimeout: time.Minute,
		ReadTimeout: time.Minute,
		WriteTimeout: time.Minute,

	})
}