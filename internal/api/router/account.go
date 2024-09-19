package router

import (
	"github.com/arifai/zenith/internal/handler"
	"github.com/arifai/zenith/internal/middleware"
	"github.com/gin-gonic/gin"
)

// AccountRouter sets up routes for account operations, including registration, authorization, and current account info fetching.
func AccountRouter(group *gin.RouterGroup, accountHandler *handler.AccountHandler, strictAuthMiddleware *middleware.StrictAuthMiddleware) {
	accountAuthGroup := group.Group("/auth/account")
	accountGroup := group.Group("/account", strictAuthMiddleware.StrictAuth())
	meGroup := accountGroup.Group("/me")
	updateGroup := meGroup.Group("/update")

	setupAccountAuthRoutes := func(g *gin.RouterGroup) {
		accountAuthGroup.POST("/registration", accountHandler.Register)
		accountAuthGroup.POST("/authorization", accountHandler.Authorization)
		accountAuthGroup.POST("/unauthorization", strictAuthMiddleware.StrictAuth(), accountHandler.Unauthorization)
	}

	setupAccountRoutes := func(g *gin.RouterGroup) {
		meGroup.GET("", accountHandler.GetCurrent)
		updateGroup.PUT("", accountHandler.Update)
		updateGroup.PUT("/password", accountHandler.UpdatePassword)
	}

	setupAccountAuthRoutes(accountAuthGroup)
	setupAccountRoutes(accountGroup)
}
