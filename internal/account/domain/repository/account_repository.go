package repository

import (
	"context"
	"github.com/arifai/go-modular-monolithic/config"
	"github.com/arifai/go-modular-monolithic/internal/account/api/types"
	"github.com/arifai/go-modular-monolithic/internal/account/domain/model"
	"github.com/arifai/go-modular-monolithic/internal/errors"
	crp "github.com/arifai/go-modular-monolithic/pkg/crypto"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AccountRepository is a struct that represent the account repository
type AccountRepository struct{ db *gorm.DB }

// NewAccountRepository creates a new account repository
func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{db: db.WithContext(context.Background()).Debug()}
}

// CreateAccount creates a new account
func (repo *AccountRepository) CreateAccount(payload *types.CreateAccountRequest) (*model.Account, error) {
	m := new(model.Account)
	exists, err := m.EmailExists(repo.db, payload.Email)
	if err != nil {
		return nil, err
	} else if exists {
		return nil, errors.ErrEmailAlreadyExists
	}

	hash := crp.Argon2IdHash{Time: 3, Memory: 64 * 1024, Threads: 4, KeyLen: 32, SaltLen: 32}
	generateHash, err := hash.GenerateHash([]byte(payload.Password), []byte(config.Load().PasswordSalt))
	if err != nil {
		return nil, err
	}
	account := &model.Account{FullName: payload.FullName, Email: payload.Email}
	accountPassHashed := &model.AccountPassHashed{
		AccountId:  account.ID,
		PassHashed: generateHash,
	}
	account.AccountPassHashed = accountPassHashed

	return account.CreateAccount(repo.db)
}

// Find finds an account by its id
func (repo *AccountRepository) Find(id uuid.UUID) (*model.Account, error) {
	account := new(model.Account)
	founded, err := account.FindByID(repo.db, id)
	if err != nil {
		return nil, err
	}

	return founded, nil
}

// FindByEmail finds an account by its email
func (repo *AccountRepository) FindByEmail(email string) (*model.Account, error) {
	account := new(model.Account)
	founded, err := account.FindByEmail(repo.db, email)
	if err != nil {
		return nil, err
	}

	return founded, nil
}
