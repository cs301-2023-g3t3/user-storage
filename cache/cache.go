package cache

import (
	"log"
	"os"

	"github.com/go-redis/redis"
)

var RedisClient *redis.Client

func ConnectToRedis() {
	var addr string
	if os.Getenv("ENV") != "lambda" {
		addr = "redis:6379"
	} else {
		addr = os.Getenv("REDIS_HOST")
		log.Println(addr)
	}
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
	_, err := RedisClient.Ping().Result()
	if err != nil {
		log.Println(err)
		log.Fatal("Failed to connect to Redis Cache")
	} else {
		log.Println("Connected to Redis Cache")
	}
}
