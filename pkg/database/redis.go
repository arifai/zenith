package database

import (
	"context"
	"fmt"
	"github.com/arifai/zenith/config"
	"github.com/redis/go-redis/v9"
	"log"
)

// ConnectRedis establishes a connection to a Redis server using the provided configuration and returns the client instance.
func ConnectRedis(config config.RedisConfig) *redis.Client {
	address := fmt.Sprintf("%s:%d", config.Host, config.Port)

	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		DB:       config.Database,
		Username: config.Username,
		Password: config.Password,
	})

	checkRedisConnection(rdb)

	return rdb
}

// checkRedisConnection pings the Redis server and logs a fatal error if the connection fails.
func checkRedisConnection(rdb *redis.Client) {
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
	}
}
