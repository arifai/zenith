package router

import (
	"github.com/arifai/zenith/config"
	"github.com/arifai/zenith/internal/account/api/handler"
	"github.com/arifai/zenith/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AccountRouter sets up all routes related to account operations under the provided gin.RouterGroup.
func AccountRouter(group *gin.RouterGroup, db *gorm.DB, config *config.Config) {
	accountHandler := handler.NewAccountHandler(db, config)

	authAccountGroup := group.Group("/auth/account")
	{
		authAccountGroup.POST("/authorization", accountHandler.AuthHandler)
		authAccountGroup.POST("/registration", accountHandler.RegisterAccountHandler)
	}

	accountGroup := group.Group("/account")
	accountGroup.Use(middleware.Middleware(db))
	{
		accountGroup.GET("/me", accountHandler.GetAccountHandler)
		accountGroup.PATCH("/me/update", accountHandler.UpdateAccountHandler)
		accountGroup.PUT("/me/update_password", accountHandler.UpdatePasswordAccountHandler)
	}
}
