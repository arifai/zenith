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

// AccountRepository handles CRUD operations for account data in the database.
type AccountRepository struct{ db *gorm.DB }

// NewAccountRepository initializes a new AccountRepository with a database context and debug mode enabled.
func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{db: db.WithContext(context.Background()).Debug()}
}

// CreateAccount registers a new user account in the system using the provided payload data. The payload must contain
// full name, email, and password. If the email already exists in the database, it returns an errors.ErrEmailAlreadyExists error.
// Passwords are hashed using Argon2ID hashing algorithm before saving. Returns the created model.Account or any errors encountered.
func (repo *AccountRepository) CreateAccount(payload *types.AccountCreateRequest) (*model.Account, error) {
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

// Find retrieves an account by its uuid.UUID from the database.
func (repo *AccountRepository) Find(id uuid.UUID) (*model.Account, error) {
	account := new(model.Account)
	founded, err := account.FindByID(repo.db, id)
	if err != nil {
		return nil, err
	}

	return founded, nil
}

// FindByEmail retrieves an account by its email address from the database.
// It returns a pointer to the model.Account and an error if any.
func (repo *AccountRepository) FindByEmail(email string) (*model.Account, error) {
	account := new(model.Account)
	founded, err := account.FindByEmail(repo.db, email)
	if err != nil {
		return nil, err
	}

	return founded, nil
}

// Update modifies an existing account with the given id using the provided payload data.
// If the account is found, its FullName and Email fields are updated and changes are saved in the database.
// Returns the updated model.Account or an error if any step fails.
func (repo *AccountRepository) Update(id uuid.UUID, payload *types.AccountUpdateRequest) (*model.Account, error) {
	founded, err := repo.Find(id)
	if err != nil {
		return nil, err
	}

	founded.FullName = payload.FullName
	founded.Email = payload.Email

	return founded.Update(repo.db)
}
