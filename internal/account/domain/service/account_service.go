package service

import (
	"fmt"
	"github.com/arifai/zenith/config"
	"github.com/arifai/zenith/internal/account/api/types"
	"github.com/arifai/zenith/internal/account/domain/model"
	"github.com/arifai/zenith/internal/account/domain/repository"
	"github.com/arifai/zenith/pkg/core"
	"github.com/arifai/zenith/pkg/errormessage"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// AccountService provides methods to manage user accounts.
type AccountService struct {
	cfg         *config.Config
	redisClient *redis.Client
	repo        *repository.AccountRepository
}

// NewAccountService initializes a new AccountService with the given database context and configuration settings.
func NewAccountService(db *gorm.DB, cfg *config.Config, redisClient *redis.Client) *AccountService {
	return &AccountService{
		cfg:         cfg,
		redisClient: redisClient,
		repo:        repository.NewAccountRepository(db, redisClient),
	}
}

// CreateAccount registers a new user account in the system using the provided payload data.
func (s *AccountService) CreateAccount(payload *types.AccountCreateRequest, config *config.Config) (*model.Account, error) {
	return s.repo.CreateAccount(payload, config)
}

// GetAccount retrieves the current account from the given context.
func (s *AccountService) GetAccount(ctx *core.Context) (*model.Account, error) {
	return s.getCurrentAccount(ctx)
}

// UpdateAccount updates the details of the current account using the provided payload.
func (s *AccountService) UpdateAccount(ctx *core.Context, payload *types.AccountUpdateRequest) (*model.Account, error) {
	currentAccount, err := s.getCurrentAccount(ctx)
	if err != nil {
		return nil, err
	}

	return s.repo.Update(currentAccount.ID, payload)
}

// UpdatePassword changes the password for the current account using the provided payload.
func (s *AccountService) UpdatePassword(ctx *core.Context, payload *types.AccountUpdatePasswordRequest, config *config.Config) (*model.Account, error) {
	currentAccount, err := s.getCurrentAccount(ctx)
	if err != nil {
		return nil, err
	}

	return s.repo.UpdatePassword(currentAccount.ID, payload, config)
}

// getCurrentAccount retrieves and asserts the current account from the context.
func (s *AccountService) getCurrentAccount(ctx *core.Context) (*model.Account, error) {
	currentAccount, ok := ctx.CurrentAccount.(*model.Account)
	if !ok {
		err := fmt.Errorf(errormessage.ErrTypeAssertionFailedText, ctx.CurrentAccount)
		return nil, err
	}

	return currentAccount, nil
}
