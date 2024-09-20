package repository

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Repository represents a repository that interacts with a database and a Redis cache.
type Repository struct {
	db    *gorm.DB
	redis *redis.Client
}

// New initializes a new Repository with the provided database and redis clients.
func New(db *gorm.DB, redis *redis.Client) *Repository {
	return &Repository{db: db, redis: redis}
}
