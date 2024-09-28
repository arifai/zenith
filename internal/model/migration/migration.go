package migration

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Migration is a wrapper around *gorm.DB, providing methods for database migrations.
type Migration struct {
	db *gorm.DB
	id uuid.UUID
}

// New initializes a new instance of Migration with the provided *gorm.DB connection.
func New(db *gorm.DB, id uuid.UUID) *Migration {
	return &Migration{db: db, id: id}
}
