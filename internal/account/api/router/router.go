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
	group.Group("/auth").
		Group("/account").
		POST("/authorization", func(c *gin.Context) {
			handler.AuthHandler(c, db, config)
		}).
		POST("/registration", func(c *gin.Context) {
			handler.RegisterAccountHandler(c, db, config)
		})

	group.Group("/account").
		Use(middleware.Middleware(db)).
		GET("/me", func(c *gin.Context) {
			handler.GetAccountHandler(c, db, config)
		}).
		PATCH("/me/update", func(c *gin.Context) {
			handler.UpdateAccountHandler(c, db, config)
		})
}
