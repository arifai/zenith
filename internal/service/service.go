package service

import (
	"github.com/arifai/zenith/config"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Service is a structure that holds the database, Redis client, and configuration settings for the service.
type Service struct {
	db     *gorm.DB
	redis  *redis.Client
	config *config.Config
}

// New initializes a new Service instance with the provided database connection, Redis client, and configuration settings.
func New(db *gorm.DB, redis *redis.Client, config *config.Config) *Service {
	return &Service{db: db, redis: redis, config: config}
}
