package router

import (
	"github.com/arifai/zenith/internal/handler"
	"github.com/arifai/zenith/internal/middleware"
	"github.com/gin-gonic/gin"
)

func NotificationRouter(group *gin.RouterGroup, notificationHandler *handler.NotificationHandler, middleware *middleware.StrictAuthMiddleware) {
	notificationGroup := group.Group("/notification", middleware.StrictAuth())
	setupNotificationRoutes := func(group *gin.RouterGroup) {
		notificationGroup.GET("/list", notificationHandler.GetList)
		notificationGroup.POST("/mark_as_read", notificationHandler.MarkAsRead)
	}

	setupNotificationRoutes(notificationGroup)
}
