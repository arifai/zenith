package server

import (
	"context"
	"fmt"
	"github.com/arifai/zenith/cmd/wire"
	cfg "github.com/arifai/zenith/cmd/wire/config"
	"github.com/arifai/zenith/cmd/wire/logger"
	"github.com/arifai/zenith/cmd/wire/migration"
	"github.com/arifai/zenith/config"
	"github.com/arifai/zenith/pkg/database"
	"github.com/arifai/zenith/pkg/errormessage"
	"github.com/arifai/zenith/pkg/tracer"
	"github.com/arifai/zenith/pkg/utils"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"os"
)

const (
	serverAddress    = ":8080"
	trustedProxyAddr = "127.0.0.1"
)

var log = logger.ProvideLogger()

// Run initializes the environment and starts the server, logging errors if the server initialization fails.
func Run() {
	fmt.Println(banner())
	log.Info("Starting server")

	tp, err := tracer.InitTracer()
	if err != nil {
		log.Error("Failed to initialize tracer", zap.Error(err))
		return
	}

	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Error("Failed to shutdown tracer provider", zap.Error(err))
		}
	}()

	initializeConfig := cfg.ProvideConfig()
	if err := initializeAndRunServer(initializeConfig); err != nil {
		log.Error(errormessage.ErrInitializingServerText, zap.Error(err))
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
func initializeAndRunServer(config *config.Config) error {
	db, err := connectDatabase(config)
	if err != nil {
		return err
	}
	rdb, err := connectRedis(config)
	if err != nil {
		return err
	}

	if err := setupRouter(db, rdb, config); err != nil {
		return err
	}

	return nil
}

func connectDatabase(config *config.Config) (*gorm.DB, error) {
	db := database.ConnectDatabase(config)
	if db == nil {
		return nil, fmt.Errorf(errormessage.ErrFailedToConnectDBText)
	}
	return db, nil
}

func connectRedis(config *config.Config) (*redis.Client, error) {
	rdb := database.ConnectRedis(config)
	if rdb == nil {
		return nil, fmt.Errorf(errormessage.ErrFailedToConnectRedisText)
	}
	return rdb, nil
}

func setupRouter(db *gorm.DB, rdb *redis.Client, config *config.Config) error {
	migrate(db)
	utils.SetupTranslation()
	rtr := wire.InitializeRouter(db, rdb, config, log)

	if err := rtr.SetTrustedProxies([]string{trustedProxyAddr}); err != nil {
		return fmt.Errorf(errormessage.ErrFailedSetTrustedProxiesText+"%v", err)
	}

	if err := rtr.Run(serverAddress); err != nil {
		return err
	}
	return nil
}

func migrate(db *gorm.DB) {
	migrator := migration.ProvideMigration(db, uuid.New(), log)
	migrator.AccountMigration()
	migrator.NotificationMigration()
}
