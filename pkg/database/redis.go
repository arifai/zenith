package database

import (
	"context"
	"fmt"
	"github.com/arifai/zenith/config"
	"github.com/arifai/zenith/pkg/errormessage"
	logg "github.com/arifai/zenith/pkg/logger"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// ConnectRedis establishes a connection to a Redis server using the provided configuration and returns the client instance.
func ConnectRedis(config *config.Config) *redis.Client {
	address := fmt.Sprintf("%s:%d", config.RedisHost, config.RedisPort)

	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		DB:       config.RedisDatabase,
		Username: config.RedisUsername,
		Password: config.RedisPassword,
	})

	checkRedisConnection(rdb)

	return rdb
}

// checkRedisConnection pings the Redis server and logs a fatal error if the connection fails.
func checkRedisConnection(rdb *redis.Client) {
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		logg.Logger.Fatal(errormessage.ErrFailedToConnectRedisText, zap.Error(err))
	}
}
