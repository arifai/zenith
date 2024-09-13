package main

import (
	"github.com/arifai/go-modular-monolithic/config"
	"github.com/arifai/go-modular-monolithic/internal/account/domain/model"
	"github.com/arifai/go-modular-monolithic/pkg/api"
	"github.com/arifai/go-modular-monolithic/pkg/database"
	"github.com/arifai/go-modular-monolithic/pkg/utils"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	cfg := config.Load()
	db := database.ConnectDatabase(cfg)
	router := gin.Default()
	model.AccountMigration(db)
	api.SetupRouter(router, db, &cfg)
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
