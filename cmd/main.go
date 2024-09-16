package main

import (
	"github.com/arifai/zenith/config"
	"github.com/arifai/zenith/internal/account/domain/model"
	"github.com/arifai/zenith/pkg/api"
	"github.com/arifai/zenith/pkg/database"
	"github.com/arifai/zenith/pkg/utils"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	envLoader := config.NewEnv(config.Config{}, config.SMTPConfig{}, config.RedisConfig{})
	defConfig := envLoader.LoadDefault()
	db := database.ConnectDatabase(defConfig)
	rdb := database.ConnectRedis(envLoader.LoadRedis())
	router := gin.Default()
	api.SetupRouter(router, db, &defConfig, rdb)
	model.AccountMigration(db)
	utils.SetupTranslation()

	trustedProxiesErr := router.SetTrustedProxies([]string{"127.0.0.1"})
	if trustedProxiesErr != nil {
		return
	}

	err := router.Run(":8080")
	if err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
