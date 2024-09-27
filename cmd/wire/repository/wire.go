//go:build wireinject

package repository

import (
	"github.com/arifai/zenith/internal/repository"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func ProvideRepository(db *gorm.DB, rdb *redis.Client) *repository.Repository {
	wire.Build(repository.New)
	return &repository.Repository{}
}

func ProvideAccountRepository(db *gorm.DB, rdb *redis.Client) repository.AccountRepository {
	wire.Build(repository.New, repository.NewAccountRepository)
	return nil
}

func ProvideNotificationRepository(db *gorm.DB, rdb *redis.Client) repository.NotificationRepository {
	wire.Build(repository.New, repository.NewNotificationRepository)
	return nil
}
