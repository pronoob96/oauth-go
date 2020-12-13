package redis

import (
	"oauth/config"
	"os"

	"github.com/go-redis/redis"
)

var Client *redis.Client

func InitializeRedis(redisConf *config.Redis) {
	dsn := os.Getenv("REDIS_DSN")
	if len(dsn) == 0 {
		dsn = redisConf.ConnectionURL
	}
	Client = redis.NewClient(&redis.Options{
		Addr: dsn,
	})
	_, err := Client.Ping().Result()
	if err != nil {
		panic(err)
	}
}
