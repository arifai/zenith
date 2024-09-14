package database

import (
	"fmt"
	"github.com/arifai/go-modular-monolithic/config"
	"github.com/arifai/go-modular-monolithic/pkg/errormessage"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

// ConnectDatabase establishes a connection to the database using the provided configuration settings.
// It returns a *gorm.DB instance for interacting with the database.
func ConnectDatabase(config config.Config) *gorm.DB {
	log.Println("Connecting to database...")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", config.DatabaseHost, config.DatabaseUser, config.DatabasePassword, config.DatabaseName, config.DatabasePort, config.SslMode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			ti, _ := time.LoadLocation(config.Timezone)

			return time.Now().In(ti)
		},
	})
	if err != nil {
		log.Fatalf("%s: %v", errormessage.ErrFailedToConnectDBText, err)
	}

	sqlDb, dbErr := db.DB()
	if dbErr != nil {
		log.Fatalf("%s: %v", errormessage.ErrFailedGetDBInstanceText, dbErr)
	}

	sqlDb.SetMaxIdleConns(10)
	sqlDb.SetMaxIdleConns(100)
	sqlDb.SetConnMaxLifetime(1 * time.Hour)

	return db
}
