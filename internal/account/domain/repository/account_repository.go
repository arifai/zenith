package repository

import (
	"context"
	"github.com/arifai/zenith/config"
	"github.com/arifai/zenith/internal/account/api/types"
	"github.com/arifai/zenith/internal/account/domain/model"
	"github.com/arifai/zenith/pkg/crypto"
	"github.com/arifai/zenith/pkg/errormessage"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AccountRepository handles CRUD operations for account data in the database.
type AccountRepository struct {
	db *gorm.DB
}

// NewAccountRepository initializes a new AccountRepository with a database context and debug mode enabled.
func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{db: db.WithContext(context.Background()).Debug()}
}

// CreateAccount registers a new user account in the system using the provided payload data.
func (repo *AccountRepository) CreateAccount(payload *types.AccountCreateRequest, config *config.Config) (*model.Account, error) {
	m := new(model.Account)
	exists, err := m.EmailExists(repo.db, payload.Email)
	if err != nil {
		return nil, err
	} else if exists {
		return nil, errormessage.ErrEmailAlreadyExists
	}

	hash, err := generatePasswordHash(payload.Password, config.PasswordSalt)
	if err != nil {
		return nil, err
	}

	account := &model.Account{FullName: payload.FullName, Email: payload.Email}
	account.AccountPassHashed = &model.AccountPassHashed{
		AccountId:  account.ID,
		PassHashed: hash,
	}

	return account.CreateAccount(repo.db)
}

func generatePasswordHash(password, salt string) (string, error) {
	return crypto.DefaultArgon2IDHash.GenerateHash([]byte(password), []byte(salt))
}

// Find retrieves an account by its uuid.UUID from the database.
func (repo *AccountRepository) Find(id uuid.UUID) (*model.Account, error) {
	account := new(model.Account)
	foundAccount, err := account.FindByID(repo.db, id)
	if err != nil {
		return nil, err
	}
	return foundAccount, nil
}

// FindByEmail retrieves an account by its email address from the database.
func (repo *AccountRepository) FindByEmail(email string) (*model.Account, error) {
	account := new(model.Account)
	foundAccount, err := account.FindByEmail(repo.db, email)
	if err != nil {
		return nil, err
	}
	return foundAccount, nil
}

// Update modifies an existing account with the given id using the provided payload data.
func (repo *AccountRepository) Update(id uuid.UUID, payload *types.AccountUpdateRequest) (*model.Account, error) {
	account, err := repo.Find(id)
	if err != nil {
		return nil, err
	}

	account.FullName = payload.FullName
	account.Email = payload.Email

	return account.Update(repo.db)
}

// UpdatePassword updates the password for an existing account identified by id.
func (repo *AccountRepository) UpdatePassword(id uuid.UUID, payload *types.AccountUpdatePasswordRequest, config *config.Config) (*model.Account, error) {
	account, err := repo.Find(id)
	if err != nil {
		return nil, err
	}

	isMatch, err := crypto.VerifyHash(payload.OldPassword, account.AccountPassHashed.PassHashed)
	if err != nil || !isMatch {
		return nil, errormessage.ErrWrongOldPassword
	}

	hash, err := generatePasswordHash(payload.NewPassword, config.PasswordSalt)
	if err != nil {
		return nil, err
	}

	account.AccountPassHashed.PassHashed = hash

	return account.UpdatePassword(repo.db)
}
