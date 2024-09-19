//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/arifai/zenith/config"
	"github.com/arifai/zenith/internal/handler"
	"github.com/arifai/zenith/internal/middleware"
	"github.com/arifai/zenith/internal/repository"
	"github.com/arifai/zenith/internal/service"
	"github.com/arifai/zenith/pkg/api"
	"github.com/arifai/zenith/pkg/common"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// InitializeResponse is a provider function that sets up and returns a new instance of common.Response using dependency injection.
func InitializeResponse() *common.Response {
	wire.Build(common.NewResponse)
	return &common.Response{}
}

// InitializeRepository sets up a new instance of repository.Repository using the provided database and redis clients.
func InitializeRepository(db *gorm.DB, redis *redis.Client) *repository.Repository {
	wire.Build(repository.New)
	return &repository.Repository{}
}

// InitializeAccountRepository initializes and returns an AccountRepository implementation.
// It uses the provided gorm.DB and redis.Client for database and caching operations.
func InitializeAccountRepository(db *gorm.DB, redis *redis.Client) repository.AccountRepository {
	wire.Build(repository.New, repository.NewAccountRepository)
	return nil
}

// InitializeService initializes a new service instance using the provided database connection, Redis client, and configuration settings.
func InitializeService(db *gorm.DB, redis *redis.Client, cfg *config.Config) *service.Service {
	wire.Build(service.New)
	return &service.Service{}
}

// InitializeAccountService sets up the AccountService by wiring dependencies including database, redis, and configuration.
func InitializeAccountService(db *gorm.DB, redis *redis.Client, cfg *config.Config) service.AccountService {
	wire.Build(service.New, repository.New, repository.NewAccountRepository, service.NewAccountService)
	return nil
}

// InitializeHandler sets up and returns a new handler.Handler instance using dependency injection.
func InitializeHandler() *handler.Handler {
	wire.Build(InitializeResponse, handler.New)
	return &handler.Handler{}
}

// InitializeAccountHandler initializes and returns a new AccountHandler by setting up necessary dependencies using wire.
func InitializeAccountHandler(db *gorm.DB, redis *redis.Client, cfg *config.Config) *handler.AccountHandler {
	wire.Build(InitializeResponse, InitializeRepository, handler.New, service.New, repository.NewAccountRepository, service.NewAccountService, handler.NewAccountHandler)
	return &handler.AccountHandler{}
}

// InitializeMiddleware initializes and returns a new middleware with the provided database and Redis client.
func InitializeMiddleware(db *gorm.DB, redis *redis.Client) *middleware.Middleware {
	wire.Build(middleware.New)
	return &middleware.Middleware{}
}

// InitializeStrictAuthMiddleware sets up and returns a new instance of StrictAuthMiddleware using the provided database and Redis clients.
func InitializeStrictAuthMiddleware(db *gorm.DB, redis *redis.Client) *middleware.StrictAuthMiddleware {
	wire.Build(middleware.New, middleware.NewStrictAuthMiddleware)
	return &middleware.StrictAuthMiddleware{}
}

// InitializeGinEngine sets up and returns a new Gin engine with routes and middleware for account handling.
func InitializeGinEngine(accountHandler *handler.AccountHandler, strictAuthMiddleware *middleware.StrictAuthMiddleware) *gin.Engine {
	engine := gin.Default()
	api.SetupRouter(engine, accountHandler, strictAuthMiddleware)
	return engine
}

// InitializeRouter sets up the Gin router with handlers and middleware for the account service.
func InitializeRouter(db *gorm.DB, redis *redis.Client, cfg *config.Config) *gin.Engine {
	wire.Build(
		InitializeAccountHandler,
		InitializeStrictAuthMiddleware,
		InitializeGinEngine,
	)
	return &gin.Engine{}
}
