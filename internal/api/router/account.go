package router

import (
	"github.com/arifai/zenith/internal/handler"
	"github.com/arifai/zenith/internal/middleware"
	"github.com/gin-gonic/gin"
)

// AccountRouter sets up routes for account operations, including registration, authorization, and current account info fetching.
func AccountRouter(group *gin.RouterGroup, accountHandler *handler.AccountHandler, middleware *middleware.StrictAuthMiddleware) {
	accountAuthGroup := group.Group("/auth/account")
	accountGroup := group.Group("/account", middleware.StrictAuth())
	meGroup := accountGroup.Group("/me")

	setupAccountAuthRoutes(accountAuthGroup, accountHandler, middleware)
	setupAccountRoutes(meGroup, accountHandler)
}

func setupAccountAuthRoutes(g *gin.RouterGroup, accountHandler *handler.AccountHandler, middleware *middleware.StrictAuthMiddleware) {
	g.POST("/registration", accountHandler.Register)
	g.POST("/authorization", accountHandler.Authorization)
	g.POST("/refresh", accountHandler.RefreshToken)
	g.POST("/unauthorization", middleware.StrictAuth(), accountHandler.Unauthorization)
}

func setupAccountRoutes(g *gin.RouterGroup, accountHandler *handler.AccountHandler) {
	g.GET("", accountHandler.GetCurrent)
	g.PUT("/update", accountHandler.Update)
	g.Group("/update").PUT("/password", accountHandler.UpdatePassword)
}
