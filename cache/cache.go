package cache

import (
	"log"
	"os"

	"github.com/go-redis/redis"
)

var RedisClient *redis.Client

func ConnectToRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"), 
		Password: "",               
		DB:       0,                
	})
	_, err := RedisClient.Ping().Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis Cache")
	} else {
		log.Println("Connected to Redis Cache")
	}
}