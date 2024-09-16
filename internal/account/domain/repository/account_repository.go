package repository

import (
	"context"
	"errors"
	"github.com/arifai/zenith/config"
	"github.com/arifai/zenith/internal/account/api/types"
	"github.com/arifai/zenith/internal/account/domain/model"
	"github.com/arifai/zenith/pkg/crypto"
	"github.com/arifai/zenith/pkg/errormessage"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// ctx is the base context used for handling timeout and cancellation signals in various operations.
var ctx = context.Background()

// AccountRepository handles CRUD operations for account data in the database.
type AccountRepository struct {
	db          *gorm.DB
	redisClient *redis.Client
}

// NewAccountRepository initializes a new AccountRepository with a database context and debug mode enabled.
func NewAccountRepository(db *gorm.DB, redisClient *redis.Client) *AccountRepository {
	return &AccountRepository{db: db.WithContext(ctx).Debug(), redisClient: redisClient}
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

// generatePasswordHash generates a secure password hash using the Argon2ID algorithm with the provided password and salt.
// Returns the password hash as a string or an error if the hashing process fails.
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

// IsTokenBlacklisted checks if a given token's jti is present in the Redis blacklist.
func (repo *AccountRepository) IsTokenBlacklisted(jti string) (bool, error) {
	value, err := repo.getRedisValue(jti)
	if errors.Is(err, redis.Nil) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return value == "blacklisted", nil
}

// getRedisValue retrieves a value from Redis using the provided key.
func (repo *AccountRepository) getRedisValue(key string) (string, error) {
	result, err := repo.redisClient.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}
