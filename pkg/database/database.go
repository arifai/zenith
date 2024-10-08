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
	"net/url"
	"time"
)

// ConnectDatabase establishes a connection to the database using the provided configuration settings.
// It returns a *gorm.DB instance for interacting with the database.
func ConnectDatabase(cfg *config.Config) *gorm.DB {
	logg.Logger.Info("connecting to database")

	dsn := buildDSN(cfg)
	nowFunc := getNowFunc(cfg.Timezone)
	return connectDatabaseWithDSN(dsn, cfg.Debug, nowFunc)
}

// connectDatabaseWithDSN is a helper function to facilitate testing by separating the creation of
// the DSN from the actual connection logic.
func connectDatabaseWithDSN(dsn string, debug bool, nowFunc func() time.Time) *gorm.DB {
	logLevel := logger.Silent
	if debug {
		logg.Logger.Info("database debug mode enabled")
		logLevel = logger.Info
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:  logger.Default.LogMode(logLevel),
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

// buildDSN constructs the Data Source Name (DSN) for connecting to a PostgreSQL database using the provided configuration settings.
func buildDSN(cfg *config.Config) string {
	query := url.Values{}
	query.Add("sslmode", cfg.SslMode)

	u := &url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.DatabaseUser, cfg.DatabasePassword),
		Host:     fmt.Sprintf("%s:%s", cfg.DatabaseHost, cfg.DatabasePort),
		Path:     cfg.DatabaseName,
		RawQuery: query.Encode(),
	}

	return u.String()
}

// getNowFunc returns a function that provides the current time in the configured timezone.
func getNowFunc(timezone string) func() time.Time {
	return func() time.Time {
		loc, _ := time.LoadLocation(timezone)
		return time.Now().In(loc)
	}
}
