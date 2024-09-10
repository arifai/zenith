package model

import (
	"gorm.io/gorm"
	"time"
)

// Account is a struct that represent the account model
type Account struct {
	ID             int64           `json:"id" gorm:"primaryKey;autoIncrement;"`
	FullName       string          `json:"full_name" gorm:"not null;column:full_name;type:varchar"`
	Email          string          `json:"email" gorm:"not null;column:email;type:varchar;uniqueIndex:idx_account_email"`
	Avatar         string          `json:"avatar" gorm:"column:avatar;type:varchar"`
	Active         bool            `json:"active" gorm:"column:active;type:boolean;default:false"`
	FcmToken       string          `json:"fcm_token" gorm:"column:fcm_token;type:varchar"`
	CreatedAt      time.Time       `json:"created_at" gorm:"column:created_at;autoCreateTime;default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time       `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	UserPassHashed *UserPassHashed `gorm:"foreignKey:AccountId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// UserPassHashed is a struct that represent the user_pass_hashed model
type UserPassHashed struct {
	ID         int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	AccountId  int64     `json:"account_id" gorm:"column:account_id;type:bigint"`
	PassHashed string    `json:"pass_hashed" gorm:"not null;column:pass_hashed;type:varchar"`
	CreatedAt  time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime;default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

// CreateAccount will create a new account
func (account *Account) CreateAccount(db *gorm.DB) (*Account, error) {
	err := db.FirstOrCreate(&account, Account{Email: account.Email}).Error
	if err := db.Model(&account).Association("UserPassHashed").Find(&account.UserPassHashed); err != nil {
		return nil, err
	}

	return account, err
}
