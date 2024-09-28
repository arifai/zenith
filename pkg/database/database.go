package database

import (
	"fmt"
	"github.com/arifai/zenith/config"
	"github.com/arifai/zenith/pkg/errormessage"
	logg "github.com/arifai/zenith/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

// ConnectDatabase establishes a connection to the database using the provided configuration settings.
// It returns a *gorm.DB instance for interacting with the database.
func ConnectDatabase(cfg *config.Config) *gorm.DB {
	logg.Logger.Info("connecting to database")

	dsn := buildDSN(cfg)
	nowFunc := getNowFunc(cfg.Timezone)
	return connectDatabaseWithDSN(dsn, nowFunc)
}

// connectDatabaseWithDSN is a helper function to facilitate testing by separating the creation of
// the DSN from the actual connection logic.
func connectDatabaseWithDSN(dsn string, nowFunc func() time.Time) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:  logger.Default.LogMode(logger.Silent),
		NowFunc: nowFunc,
	})
	if err != nil {
		logg.Logger.Fatal(errormessage.ErrFailedToConnectDBText, zap.Error(err))
	}

	sqlDb, sqlDBError := db.DB()
	if sqlDBError != nil {
		logg.Logger.Fatal(errormessage.ErrFailedGetDBInstanceText, zap.Error(sqlDBError))
	}

	sqlDb.SetMaxIdleConns(10)
	sqlDb.SetMaxOpenConns(100)
	sqlDb.SetConnMaxLifetime(1 * time.Hour)

	return db
}

// buildDSN constructs the database source name from the provided configuration.
func buildDSN(cfg *config.Config) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.DatabaseHost, cfg.DatabaseUser, cfg.DatabasePassword, cfg.DatabaseName,
		cfg.DatabasePort, cfg.SslMode)
}

// getNowFunc returns a function that provides the current time in the configured timezone.
func getNowFunc(timezone string) func() time.Time {
	return func() time.Time {
		loc, _ := time.LoadLocation(timezone)
		return time.Now().In(loc)
	}
}
