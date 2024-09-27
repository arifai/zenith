//go:build wireinject

package handler

import (
	cmn "github.com/arifai/zenith/cmd/wire/common"
	repo "github.com/arifai/zenith/cmd/wire/repository"
	"github.com/arifai/zenith/config"
	"github.com/arifai/zenith/internal/handler"
	"github.com/arifai/zenith/internal/repository"
	"github.com/arifai/zenith/internal/service"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func ProvideHandler() *handler.Handler {
	wire.Build(cmn.ProvideResponse, handler.New)
	return &handler.Handler{}
}

func ProvideAccountHandler(db *gorm.DB, rdb *redis.Client, cfg *config.Config) *handler.AccountHandler {
	wire.Build(cmn.ProvideResponse, repo.ProvideRepository, handler.New, service.New, repository.NewAccountRepository, service.NewAccountService, handler.NewAccountHandler)
	return &handler.AccountHandler{}
}

func ProvideNotificationHandler(db *gorm.DB, rdb *redis.Client, cfg *config.Config) *handler.NotificationHandler {
	wire.Build(cmn.ProvideResponse, repo.ProvideRepository, handler.New, service.New, repository.NewNotificationRepository, service.NewNotificationService, handler.NewNotificationHandler)
	return &handler.NotificationHandler{}
}
