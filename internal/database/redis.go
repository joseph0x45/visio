package database

import (
	"fmt"
	"os"
	"github.com/redis/go-redis/v9"
)

func GetRedisClient() *redis.Client {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		fmt.Println("Failed to read REDIS_URL environment variable")
		os.Exit(1)
	}
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		fmt.Printf("Error while parsing redis URL: %v", err)
		os.Exit(1)
	}
	return redis.NewClient(opts)
}
