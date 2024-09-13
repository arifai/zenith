package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	AccountPassHashed *AccountPassHashed `json:"-" gorm:"foreignKey:AccountId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// AccountPassHashed is a struct that represent the user_pass_hashed model
type AccountPassHashed struct {
	ID         uuid.UUID  `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	AccountId  uuid.UUID  `json:"account_id" gorm:"not null;column:account_id;type:uuid"`
	PassHashed string     `json:"pass_hashed" gorm:"not null;column:pass_hashed;type:varchar"`
	CreatedAt  time.Time  `json:"created_at" gorm:"column:created_at;autoCreateTime;default:CURRENT_TIMESTAMP"`
	UpdatedAt  *time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
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
		if err := tx.Where(AccountPassHashed{AccountId: a.ID}).FirstOrCreate(a.AccountPassHashed).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return a, nil
}

// FindByID finds an Account by its ID
func (a *Account) FindByID(db *gorm.DB, id uuid.UUID) (*Account, error) {
	if err := db.First(a, id).Error; err != nil {
		return nil, err
	}

	return a, nil
}

// FindByEmail finds an Account by its email
func (a *Account) FindByEmail(db *gorm.DB, email string) (*Account, error) {
	if err := db.Where(&Account{Email: email}).Preload("AccountPassHashed").First(a).Error; err != nil {
		return nil, err
	}

	return a, nil
}

// EmailExists checks if an email already exists in the database
func (a *Account) EmailExists(db *gorm.DB, email string) (bool, error) {
	var count int64
	if err := db.Model(&Account{}).Where(&Account{Email: email}).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// Update updates an Account
func (a *Account) Update(db *gorm.DB) (*Account, error) {
	if err := db.Save(a).Clauses(clause.Returning{}).Error; err != nil {
		return nil, err
	}

	return a, nil
}
