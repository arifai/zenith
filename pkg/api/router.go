package api

import (
	"github.com/arifai/zenith/internal/api/router"
	"github.com/arifai/zenith/internal/handler"
	"github.com/arifai/zenith/internal/middleware"
	"github.com/gin-gonic/gin"
)

// SetupRouter initializes the main router and sets up all the routes and groups under "/api/v1".
func SetupRouter(engine *gin.Engine, accountHandler *handler.AccountHandler, notificationHandler *handler.NotificationHandler, middleware *middleware.StrictAuthMiddleware) *gin.Engine {
	apiV1 := engine.Group("/api/v1")
	router.AccountRouter(apiV1, accountHandler, middleware)
	router.NotificationRouter(apiV1, notificationHandler, middleware)
	return engine
}
