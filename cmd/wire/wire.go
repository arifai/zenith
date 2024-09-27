//go:build wireinject

package wire

import (
	"github.com/arifai/zenith/cmd/wire/handler"
	"github.com/arifai/zenith/cmd/wire/middleware"
	"github.com/arifai/zenith/config"
	"github.com/arifai/zenith/pkg/server/http"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func InitializeRouter(db *gorm.DB, redis *redis.Client, cfg *config.Config) *gin.Engine {
	wire.Build(
		handler.ProvideAccountHandler,
		handler.ProvideNotificationHandler,
		middleware.WireMiddlewareSet,
		http.ProvideGinEngine,
	)
	return &gin.Engine{}
}
