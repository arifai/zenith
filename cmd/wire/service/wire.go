//go:build wireinject

package service

import (
	"github.com/arifai/zenith/config"
	"github.com/arifai/zenith/internal/repository"
	"github.com/arifai/zenith/internal/service"
	"github.com/arifai/zenith/pkg/logger"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func ProvideService(cfg *config.Config, log logger.Logger) *service.Service {
	wire.Build(service.New)
	return &service.Service{}
}

func ProvideAccountService(db *gorm.DB, rdb *redis.Client, cfg *config.Config, log logger.Logger) service.AccountService {
	wire.Build(service.New, repository.New, repository.NewAccountRepository, service.NewAccountService)
	return nil
}

func ProvideNotificationService(db *gorm.DB, rdb *redis.Client, cfg *config.Config, log logger.Logger) service.NotificationService {
	wire.Build(service.New, repository.New, repository.NewNotificationRepository, service.NewNotificationService)
	return nil
}
