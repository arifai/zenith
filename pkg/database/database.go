package database

import (
	"fmt"
	"log"
	"time"

	"github.com/arifai/go-modular-monolithic/config"
	"github.com/arifai/go-modular-monolithic/pkg/errormessage"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ConnectDatabase establishes a connection to the database using the provided configuration settings.
// It returns a *gorm.DB instance for interacting with the database.
func ConnectDatabase(cfg config.Config) *gorm.DB {
	log.Println("Connecting to database...")

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
		log.Fatalf("%s: %v", errormessage.ErrFailedToConnectDBText, err)
	}

	sqlDb, sqlDBError := db.DB()
	if sqlDBError != nil {
		log.Fatalf("%s: %v", errormessage.ErrFailedGetDBInstanceText, sqlDBError)
	}

	sqlDb.SetMaxIdleConns(10)
	sqlDb.SetMaxOpenConns(100)
	sqlDb.SetConnMaxLifetime(1 * time.Hour)

	return db
}

// buildDSN constructs the database source name from the provided configuration.
func buildDSN(cfg config.Config) string {
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
