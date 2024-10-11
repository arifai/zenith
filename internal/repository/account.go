package repository

import (
	"context"
	"github.com/arifai/zenith/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type (
	// AccountRepository defines methods to interact with account-related data in the database.
	AccountRepository interface {
		// Create inserts a new account with the given hashed password into the database.
		Create(account *model.Account, hashedPassword string) error

		// FindByEmail retrieves an account by its email address from the database. Returns the account and any error encountered.
		FindByEmail(email string) (*model.Account, error)

		// FindByID retrieves an account by its unique identifier. Returns the account model and any error encountered.
		FindByID(id *uuid.UUID) (*model.Account, error)

		// Update updates the details of an existing account in the database. Returns an error if the operation fails.
		Update(account *model.Account) error

		// UpdatePassword updates the hashed password of an account in the database. Returns an error if the update operation fails.
		UpdatePassword(account *model.Account) error

		// SetFCMToken updates the FCM token for an account identified by the provided email address. Returns an error if the update fails.
		SetFCMToken(email, fcmToken string) error

		// UnsetFCMToken unsets the FCM token for the account associated with the specified UUID, removing the existing token if present.
		UnsetFCMToken(id uuid.UUID) error

		// BlacklistToken adds the token identified by jti to the blacklist, exp indicates the token's expiration time.
		BlacklistToken(jti string, exp time.Time) error
	}

	// accountRepository encapsulates a Repository to provide specific methods for handling account data.
	accountRepository struct{ *Repository }
)

// NewAccountRepository returns an implementation of AccountRepository using the provided Repository.
func NewAccountRepository(r *Repository) AccountRepository {
	return &accountRepository{Repository: r}
}

func (a *accountRepository) Create(account *model.Account, hashedPassword string) error {
	return a.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(account).Error; err != nil {
			return err
		}

		accountPassHashed := &model.AccountPassHashed{AccountId: account.ID, PassHashed: hashedPassword}
		if err := tx.Create(accountPassHashed).Error; err != nil {
			return err
		}

		return nil
	})
}

func (a *accountRepository) FindByEmail(email string) (*model.Account, error) {
	account := &model.Account{}
	if err := a.db.Where(&model.Account{Email: email}).Preload("AccountPassHashed").
		First(account).Error; err != nil {
		return nil, err
	}

	return account, nil
}

func (a *accountRepository) FindByID(id *uuid.UUID) (*model.Account, error) {
	var account model.Account
	if err := a.db.Preload("AccountPassHashed").Where("id = ?", id).First(&account).Error; err != nil {
		return nil, err
	}

	return &account, nil
}

func (a *accountRepository) Update(account *model.Account) error {
	return a.db.Clauses(clause.Returning{}).Save(account).Error
}

func (a *accountRepository) UpdatePassword(account *model.Account) error {
	return a.db.Model(&model.AccountPassHashed{}).
		Clauses(clause.Returning{}).
		Where("account_id = ?", account.ID).
		Update("pass_hashed", account.AccountPassHashed.PassHashed).Error
}

func (a *accountRepository) SetFCMToken(email, fcmToken string) error {
	return a.db.Model(&model.Account{}).Where(&model.Account{Email: email}).
		Update("fcm_token", fcmToken).Error
}

func (a *accountRepository) UnsetFCMToken(id uuid.UUID) error {
	return a.db.Model(&model.Account{}).Where(&model.Account{ID: id}).
		Update("fcm_token", nil).Error
}

func (a *accountRepository) BlacklistToken(jti string, exp time.Time) error {
	ttl := time.Until(exp)

	if err := a.redis.Set(context.Background(), jti, "blacklisted", ttl).Err(); err != nil {
		return err
	}

	return nil
}
