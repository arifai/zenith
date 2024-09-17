package router

import (
	"github.com/arifai/zenith/config"
	"github.com/arifai/zenith/internal/account/api/handler"
	"github.com/arifai/zenith/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// AccountRouter sets up all routes related to account operations under the provided gin.RouterGroup.
func AccountRouter(group *gin.RouterGroup, db *gorm.DB, config *config.Config, redisClient *redis.Client) {
	accountHandler := handler.NewAccountHandler(db, config, redisClient)

	authAccountGroup := group.Group("/auth/account")
	{
		authAccountGroup.POST("/registration", accountHandler.RegisterAccountHandler)
		authAccountGroup.POST("/authorization", accountHandler.AuthHandler)
		authAccountGroup.POST("/unauthorization", accountHandler.UnauthHandler).
			Use(middleware.Middleware(db, redisClient))
		authAccountGroup.POST("/refresh_token", accountHandler.RefreshTokenHandler).
			Use(middleware.Middleware(db, redisClient))
	}

	accountGroup := group.Group("/account")
	accountGroup.Use(middleware.Middleware(db, redisClient))
	{
		accountGroup.GET("/me", accountHandler.GetAccountHandler)
		accountGroup.PATCH("/me/update", accountHandler.UpdateAccountHandler)
		accountGroup.PUT("/me/update_password", accountHandler.UpdatePasswordAccountHandler)
	}
}
