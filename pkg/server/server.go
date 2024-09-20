package server

import (
	"fmt"
	"github.com/arifai/zenith/cmd/wire"
	"github.com/arifai/zenith/config"
	"github.com/arifai/zenith/internal/model/migration"
	"github.com/arifai/zenith/pkg/database"
	"github.com/arifai/zenith/pkg/utils"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log"
	"os"
)

const (
	serverAddress    = ":8080"
	trustedProxyAddr = "127.0.0.1"
)

// RunServer initializes the environment and starts the server, logging errors if the server initialization fails.
func RunServer() {
	fmt.Println(banner())
	envLoader := config.NewEnv(config.Config{}, config.SMTPConfig{}, config.RedisConfig{})
	if err := initializeAndRunServer(envLoader); err != nil {
		log.Fatalf("Error initializing server: %v", err)
	}
}

// banner reads the contents of "ascii.txt" and returns it as a string. If an error occurs, it returns the error message.
func banner() string {
	b, err := os.ReadFile("pkg/server/ascii.txt")
	if err != nil {
		return err.Error()
	}

	return string(b)
}

// initializeAndRunServer initializes the environment, connects to the database and Redis, and sets up the server router.
func initializeAndRunServer(envLoader *config.EnvImpl) error {
	defConfig := envLoader.LoadDefault()
	db, err := connectDatabase(defConfig)
	if err != nil {
		return err
	}
	rdb, err := connectRedis(envLoader)
	if err != nil {
		return err
	}

	if err := setupRouter(db, rdb, defConfig); err != nil {
		return err
	}

	return nil
}

func connectDatabase(cfg config.Config) (*gorm.DB, error) {
	db := database.ConnectDatabase(cfg)
	if db == nil {
		return nil, fmt.Errorf("failed to connect to the database")
	}
	return db, nil
}

func connectRedis(envLoader *config.EnvImpl) (*redis.Client, error) {
	rdb := database.ConnectRedis(envLoader.LoadRedis())
	if rdb == nil {
		return nil, fmt.Errorf("failed to connect to Redis")
	}
	return rdb, nil
}

func setupRouter(db *gorm.DB, rdb *redis.Client, defConfig config.Config) error {
	migration.AccountMigration(db)
	utils.SetupTranslation()

	rtr := wire.InitializeRouter(db, rdb, &defConfig)

	if err := rtr.SetTrustedProxies([]string{trustedProxyAddr}); err != nil {
		return fmt.Errorf("failed to set trusted proxies: %v", err)
	}

	if err := rtr.Run(serverAddress); err != nil {
		return err
	}
	return nil
}
