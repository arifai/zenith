package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// Account is a struct that represent the account model
type Account struct {
	ID                uuid.UUID          `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	FullName          string             `json:"full_name" gorm:"not null;column:full_name;type:varchar"`
	Email             string             `json:"email" gorm:"not null;column:email;type:varchar;uniqueIndex:idx_account_email"`
	Avatar            string             `json:"avatar" gorm:"column:avatar;type:varchar"`
	Active            bool               `json:"active" gorm:"column:active;type:boolean;default:false"`
	FcmToken          string             `json:"fcm_token" gorm:"column:fcm_token;type:varchar"`
	CreatedAt         time.Time          `json:"created_at" gorm:"column:created_at;autoCreateTime;default:CURRENT_TIMESTAMP"`
	UpdatedAt         *time.Time         `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	AccountPassHashed *AccountPassHashed `gorm:"foreignKey:AccountId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// AccountPassHashed is a struct that represent the user_pass_hashed model
type AccountPassHashed struct {
	ID         uuid.UUID  `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	AccountId  uuid.UUID  `json:"account_id" gorm:"not null;column:account_id;type:uuid"`
	PassHashed string     `json:"pass_hashed" gorm:"not null;column:pass_hashed;type:varchar"`
	CreatedAt  time.Time  `json:"created_at" gorm:"column:created_at;autoCreateTime;default:CURRENT_TIMESTAMP"`
	UpdatedAt  *time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (a *Account) BeforeCreate() (err error) {
	a.ID = uuid.New()

	return
}

func (a *AccountPassHashed) BeforeCreate() (err error) {
	a.ID = uuid.New()

	return
}

// CreateAccount creates a new account and handles its association with AccountPassHashed
func (a *Account) CreateAccount(db *gorm.DB) (*Account, error) {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.FirstOrCreate(a, Account{Email: a.Email}).Error; err != nil {
		tx.Rollback()

		return nil, err
	}

	if a.AccountPassHashed != nil {
		a.AccountPassHashed.AccountId = a.ID
		if err := tx.Create(a.AccountPassHashed).Error; err != nil {
			tx.Rollback()

			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return a, nil
}
