package api

import (
	"github.com/arifai/zenith/config"
	"github.com/arifai/zenith/internal/account/api/router"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRouter initializes the main router and sets up all the routes and groups under "/api/v1".
func SetupRouter(engine *gin.Engine, db *gorm.DB, config *config.Config) {
	apiV1 := engine.Group("/api/v1")
	router.AccountRouter(apiV1, db, config)
}
