package model

import (
	"database/sql/driver"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type (
	Platform string

	Status string

	// PushNotification represents a system or application event message that can be presented to users.
	PushNotification struct {
		ID        uuid.UUID         `json:"id" gorm:"not null;primaryKey;type:uuid;default:uuid_generate_v4()"`
		AccountId uuid.UUID         `json:"account_id" gorm:"not null;column:account_id;type:uuid;index:idx_push_notification_account_id,hash"`
		Title     string            `json:"title" gorm:"not null;column:title;type:varchar"`
		Message   string            `json:"message" gorm:"not null;column:message;type:varchar"`
		Image     string            `json:"image" gorm:"column:image;type:varchar"`
		Data      map[string]string `json:"data" gorm:"not null;column:data;type:jsonb;serializer:json"`
		Platform  Platform          `json:"platform" gorm:"not null;column:platform;type:platform"`
		Status    Status            `json:"status" gorm:"not null;column:status;type:status;default:'Pending'"`
		Retries   int8              `json:"retries" gorm:"not null;column:retries;type:smallint;default:0"`
		CreatedAt time.Time         `json:"created_at" gorm:"column:created_at;autoCreateTime;default:CURRENT_TIMESTAMP"`
	}

	Notification struct {
		ID               uuid.UUID  `json:"id" gorm:"not null;primaryKey;type:uuid;default:uuid_generate_v4()"`
		AccountID        uuid.UUID  `json:"account_id" gorm:"not null;column:account_id;type:uuid;index:idx_notification_account_id,hash"`
		Title            string     `json:"title" gorm:"not null;column:title;type:varchar"`
		Image            string     `json:"image" gorm:"column:image;type:varchar"`
		ShortDescription string     `json:"short_description" gorm:"not null;column:short_description;type:varchar"`
		Description      string     `json:"description" gorm:"not null;column:description;type:text"`
		Read             bool       `json:"read" gorm:"not null;column:read;type:boolean;default:false"`
		ReadAt           *time.Time `json:"read_at" gorm:"column:read_at;type:timestamp"`
		CreatedAt        time.Time  `json:"created_at" gorm:"column:created_at;autoCreateTime;default:CURRENT_TIMESTAMP"`
	}
)

const (
	Android Platform = "Android"
	IOS     Platform = "iOS"
	Web     Platform = "Web"
	Pending Status   = "Pending"
	Success Status   = "Success"
	Failure Status   = "Failure"
)

func (p *Platform) Scan(value interface{}) error {
	if strValue, ok := value.(string); ok {
		*p = Platform(strValue)
		return nil
	}

	return errors.New("failed to scan platform")
}

func (p Platform) Value() (driver.Value, error) {
	return string(p), nil
}

func (s *Status) Scan(value interface{}) error {
	if strValue, ok := value.(string); ok {
		*s = Status(strValue)
		return nil
	}

	return errors.New("failed to scan status")
}

func (s Status) Value() (driver.Value, error) {
	return string(s), nil
}

func (n *PushNotification) BeforeSave(db *gorm.DB) error {
	if n.Platform != Android && n.Platform != IOS && n.Platform != Web {
		return errors.New("invalid platform")
	}

	if n.Status != Pending && n.Status != Success && n.Status != Failure {
		return errors.New("invalid status")
	}

	return nil
}
