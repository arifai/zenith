package repository

import (
	"github.com/arifai/zenith/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

type (
	// AccountRepository defines methods to interact with account-related data in the database.
	AccountRepository interface {
		// Create inserts a new account record into the database. Returns an error if the operation fails.
		Create(account *model.Account) error

		// FindByEmail retrieves an account by its email address from the database. Returns the account and any error encountered.
		FindByEmail(email string) (*model.Account, error)

		// FindByID retrieves an account by its unique identifier. Returns the account model and any error encountered.
		FindByID(id *uuid.UUID) (*model.Account, error)

		// Update updates the details of an existing account in the database. Returns an error if the operation fails.
		Update(account *model.Account) error

		// UpdatePassword updates the hashed password of an account in the database. Returns an error if the update operation fails.
		UpdatePassword(account *model.Account) error

		SetFCMToken(email, fcmToken string) error

		RemoveFCMToken(id uuid.UUID) error
	}

	// accountRepository encapsulates a Repository to provide specific methods for handling account data.
	accountRepository struct{ *Repository }
)

// NewAccountRepository returns an implementation of AccountRepository using the provided Repository.
func NewAccountRepository(r *Repository) AccountRepository {
	return &accountRepository{Repository: r}
}

func (a *accountRepository) Create(account *model.Account) error {
	return a.db.Create(account).Error
}

func (a *accountRepository) FindByEmail(email string) (*model.Account, error) {
	account := new(model.Account)
	if err := a.db.Where(&model.Account{Email: email}).Preload("AccountPassHashed").
		First(account).Error; err != nil {
		return nil, err
	}

	return account, nil
}

func (a *accountRepository) FindByID(id *uuid.UUID) (*model.Account, error) {
	account := new(model.Account)
	if err := a.db.First(account, id).Error; err != nil {
		return nil, err
	}

	return account, nil
}

func (a *accountRepository) Update(account *model.Account) error {
	return a.db.Clauses(clause.Returning{}).Save(account).Error
}

func (a *accountRepository) UpdatePassword(account *model.Account) error {
	return a.db.Where(&model.AccountPassHashed{AccountId: account.ID}).
		Update("pass_hashed", account.AccountPassHashed.PassHashed).Error
}

func (a *accountRepository) SetFCMToken(email, fcmToken string) error {
	return a.db.Model(&model.Account{}).Where(&model.Account{Email: email}).
		Update("fcm_token", fcmToken).Error
}

func (a *accountRepository) RemoveFCMToken(id uuid.UUID) error {
	return a.db.Model(&model.Account{}).Where(&model.Account{ID: id}).
		Update("fcm_token", nil).Error
}
