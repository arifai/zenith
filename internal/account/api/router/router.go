package router

import (
	"github.com/arifai/zenith/config"
	"github.com/arifai/zenith/internal/account/api/handler"
	"github.com/arifai/zenith/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// AccountRouter sets up route groups for account-related operations, including authentication and account management.
func AccountRouter(group *gin.RouterGroup, db *gorm.DB, config *config.Config, redisClient *redis.Client) {
	accountHandler := handler.NewAccountHandler(db, config, redisClient)
	middlewareFunc := middleware.Middleware(db, redisClient)

	setupAccountAuthRoutes := func(accountAuthGroup *gin.RouterGroup) {
		accountAuthGroup.POST("/registration", accountHandler.RegisterAccountHandler)
		accountAuthGroup.POST("/authorization", accountHandler.AuthHandler)
		accountAuthGroup.POST("/unauthorization", middlewareFunc, accountHandler.UnauthHandler)
		// @FIXME: `refresh_token` no need middleware
		accountAuthGroup.POST("/refresh_token", middlewareFunc, accountHandler.RefreshTokenHandler)
	}

	setupAccountRoutes := func(accountGroup *gin.RouterGroup) {
		accountGroup.Use(middlewareFunc)
		accountGroup.GET("/me", accountHandler.GetAccountHandler)
		accountGroup.PATCH("/me/update", accountHandler.UpdateAccountHandler)
		accountGroup.PUT("/me/update_password", accountHandler.UpdatePasswordAccountHandler)
	}

	authGroup := group.Group("/auth/account")
	setupAccountAuthRoutes(authGroup)

	accountGroup := group.Group("/account")
	setupAccountRoutes(accountGroup)
}
