package http

import (
	"github.com/arifai/zenith/internal/handler"
	"github.com/arifai/zenith/internal/middleware"
	"github.com/arifai/zenith/pkg/api"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func ProvideGinEngine(accountHandler *handler.AccountHandler, notificationHandler *handler.NotificationHandler, mid *middleware.StrictAuthMiddleware) *gin.Engine {
	engine := gin.Default()
	engine.Use(otelgin.Middleware("zenith-server"))
	api.SetupRouter(engine, accountHandler, notificationHandler, mid)

	return engine
}
