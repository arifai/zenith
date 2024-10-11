package migration

import (
	"github.com/arifai/zenith/pkg/logger"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Migration enables smooth database schema transitions and logs operations for debugging and auditing purposes.
type Migration struct {
	*gorm.DB
	uuid.UUID
	logger.Logger
}

// New initializes a new Migration instance with the provided database connection, UUID, and logger.
func New(db *gorm.DB, id uuid.UUID, log logger.Logger) *Migration {
	return &Migration{db, id, log}
}
