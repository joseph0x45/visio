package database

import (
	"os"
	"github.com/redis/go-redis/v9"
)

func GetRedisClient() *redis.Client {
	opts, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		panic(err)
	}
	client := redis.NewClient(opts)
	return client
}
