package migration

import "gorm.io/gorm"

// Migration is a wrapper around *gorm.DB, providing methods for database migrations.
type Migration struct{ *gorm.DB }

// New initializes a new instance of Migration with the provided *gorm.DB connection.
func New(db *gorm.DB) *Migration {
	return &Migration{db}
}
