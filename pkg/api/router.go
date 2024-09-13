package api

import (
	"github.com/arifai/go-modular-monolithic/config"
	"github.com/arifai/go-modular-monolithic/internal/account/api/router"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRouter is a function to set up the router
func SetupRouter(engine *gin.Engine, db *gorm.DB, config *config.Config) {
	apiV1 := engine.Group("/api/v1")
	router.AccountRouter(apiV1, db, config)
}
