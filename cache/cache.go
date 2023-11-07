package cache

import (
	"context"
	"crypto/tls"
	"os"

	"github.com/redis/go-redis/v9"
)


var RedisClient *redis.ClusterClient

func ConnectToRedis() {
	var node1, node2 string
	if os.Getenv("ENV") != "lambda" {
		node1 = "redis:6379"
	} else {
		node1 = os.Getenv("REDIS_NODE_1")
		node2 = os.Getenv("REDIS_NODE_2")
	}
	
	// RedisClient = redis.NewClient(&redis.Options{
	// 	TLSConfig: &tls.Config{
	// 		MinVersion: tls.VersionTLS12,
	// 	},
	// 	Addr:     addr,
	// 	Password: "",
	// 	DB:       0,
	// })

	RedisClient = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:          []string{node1, node2},
		TLSConfig:      &tls.Config{},
		ReadOnly:       false,
		RouteRandomly:  false,
		RouteByLatency: false,
	})

	ctx := context.Background()
	err := RedisClient.ForEachShard(ctx, func(ctx context.Context, shard *redis.Client) error {
		return shard.Ping(ctx).Err()
	})
	if err != nil {
		panic(err)
	}
}
