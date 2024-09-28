package middleware

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Middleware provides the necessary dependencies for handling requests, including database, Redis client, and account repository.
type Middleware struct {
	db    *gorm.DB
	redis *redis.Client
}

// New initializes and returns a new Middleware struct with the provided database, Redis client, and account repository.
func New(db *gorm.DB, redis *redis.Client) *Middleware {
	return &Middleware{
		db:    db,
		redis: redis,
	}
}
