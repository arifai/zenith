package http

import (
	"github.com/arifai/zenith/internal/handler"
	"github.com/arifai/zenith/internal/middleware"
	"github.com/arifai/zenith/pkg/api"
	"github.com/gin-gonic/gin"
)

func ProvideGinEngine(accountHandler *handler.AccountHandler, notificationHandler *handler.NotificationHandler, mid *middleware.StrictAuthMiddleware) *gin.Engine {
	engine := gin.Default()
	api.SetupRouter(engine, accountHandler, notificationHandler, mid)

	return engine
}
