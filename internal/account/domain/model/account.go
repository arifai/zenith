package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	AccountId  uuid.UUID  `json:"account_id" gorm:"not null;column:account_id;type:uuid"`
	PassHashed string     `json:"pass_hashed" gorm:"not null;column:pass_hashed;type:varchar"`
	CreatedAt  time.Time  `json:"created_at" gorm:"column:created_at;autoCreateTime;default:CURRENT_TIMESTAMP"`
	UpdatedAt  *time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

// CreateAccount creates a new Account in the database. Returns the created Account or an error.
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

// FindByID retrieves an account from the database using the provided [uuid.UUID].
// Returns the Account and an error if encountered.
func (a *Account) FindByID(db *gorm.DB, id uuid.UUID) (*Account, error) {
	if err := db.Preload("AccountPassHashed").First(a, id).Error; err != nil {
		return nil, err
	}

	return a, nil
}

// FindByEmail retrieves an account based on the provided email. It preloads the associated AccountPassHashed data.
// Returns the Account and an error if encountered.
func (a *Account) FindByEmail(db *gorm.DB, email string) (*Account, error) {
	if err := db.Where(&Account{Email: email}).Preload("AccountPassHashed").First(a).Error; err != nil {
		return nil, err
	}

	return a, nil
}

// EmailExists checks if an email already exists in the database.
// Returns true if the email exists, false otherwise, or an error.
func (a *Account) EmailExists(db *gorm.DB, email string) (bool, error) {
	var count int64
	if err := db.Where(&Account{Email: email}).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// Update updates the current Account entity in the database. Returns the updated Account or an error.
func (a *Account) Update(db *gorm.DB) (*Account, error) {
	if err := db.Updates(a).Clauses(clause.Returning{}).Error; err != nil {
		return nil, err
	}

	return a, nil
}

// UpdatePassword updates the password of an Account in the database. Returns the updated Account or an error.
func (a *Account) UpdatePassword(db *gorm.DB) (*Account, error) {
	if err := db.Where(AccountPassHashed{AccountId: a.ID}).
		Updates(AccountPassHashed{PassHashed: a.AccountPassHashed.PassHashed}).
		Error; err != nil {
		return nil, err
	}

	return a, nil
}
