package model

import (
	"github.com/google/uuid"
	"time"
)

// Account represents a user entity with various attributes such as ID, FullName, Email, etc.
// It includes fields for account management and profile information.
type Account struct {
	ID                uuid.UUID          `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	FullName          string             `json:"full_name" gorm:"not null;column:full_name;type:varchar"`
	Email             string             `json:"email" gorm:"not null;column:email;type:varchar;uniqueIndex:idx_account_email"`
	Avatar            string             `json:"avatar" gorm:"column:avatar;type:varchar"`
	Active            bool               `json:"active" gorm:"column:active;type:boolean;default:false"`
	FcmToken          string             `json:"fcm_token" gorm:"column:fcm_token;type:varchar"`
	CreatedAt         time.Time          `json:"created_at" gorm:"column:created_at;autoCreateTime;default:CURRENT_TIMESTAMP"`
	UpdatedAt         *time.Time         `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	AccountPassHashed *AccountPassHashed `json:"-" gorm:"foreignKey:AccountId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// AccountPassHashed represents a hashed password associated with an account.
type AccountPassHashed struct {
	ID         uuid.UUID  `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	AccountId  uuid.UUID  `json:"account_id" gorm:"not null;column:account_id;type:uuid;index:idx_account_pass_hashed_account_id,hash"`
	PassHashed string     `json:"pass_hashed" gorm:"not null;column:pass_hashed;type:varchar"`
	CreatedAt  time.Time  `json:"created_at" gorm:"column:created_at;autoCreateTime;default:CURRENT_TIMESTAMP"`
	UpdatedAt  *time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}
