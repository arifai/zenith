package router

import (
	"github.com/arifai/zenith/internal/handler"
	"github.com/arifai/zenith/internal/middleware"
	"github.com/gin-gonic/gin"
)

// NotificationRouter sets up routes for handling notification-related requests with required middleware.
func NotificationRouter(group *gin.RouterGroup, notificationHandler *handler.NotificationHandler, middleware *middleware.StrictAuthMiddleware) {
	notificationGroup := group.Group("/notification", middleware.StrictAuth())

	setupNotificationRoutes(notificationGroup, notificationHandler)
}

func setupNotificationRoutes(group *gin.RouterGroup, notificationHandler *handler.NotificationHandler) {
	group.GET("/list", notificationHandler.GetList)
	group.POST("/mark_as_read", notificationHandler.MarkAsRead)
}
