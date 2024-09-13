package router

import (
	"github.com/arifai/go-modular-monolithic/config"
	"github.com/arifai/go-modular-monolithic/internal/account/api/handler"
	"github.com/arifai/go-modular-monolithic/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AccountRouter is a function to handle the account router
func AccountRouter(group *gin.RouterGroup, db *gorm.DB, config *config.Config) {
	accountGroup := group.Group("/account")
	accountGroup.POST("/authorization", func(c *gin.Context) {
		handler.AuthHandler(c, db, config)
	})
	accountGroup.GET("/me", middleware.Middleware(db), func(c *gin.Context) {
		handler.GetAccountHandler(c, db, config)
	})
}
