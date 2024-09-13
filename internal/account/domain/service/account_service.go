package service

import (
	"github.com/arifai/go-modular-monolithic/config"
	"github.com/arifai/go-modular-monolithic/internal/account/api/types"
	"github.com/arifai/go-modular-monolithic/internal/account/domain/model"
	"github.com/arifai/go-modular-monolithic/internal/account/domain/repository"
	"github.com/arifai/go-modular-monolithic/pkg/core"
	"gorm.io/gorm"
	"log"
)

// AccountService is a struct that represent the account service
type AccountService struct {
	config *config.Config
	repo   *repository.AccountRepository
}

// NewAccountService creates a new account service
func NewAccountService(db *gorm.DB, config *config.Config) *AccountService {
	return &AccountService{
		config: config,
		repo:   repository.NewAccountRepository(db),
	}
}

// CreateAccount creates a new account
func (s *AccountService) CreateAccount(payload *types.CreateAccountRequest) (*model.Account, error) {
	return s.repo.CreateAccount(payload)
}

// GetAccount gets the current account
func (s *AccountService) GetAccount(ctx *core.Context) (m *model.Account, err error) {
	currentAccount, ok := ctx.CurrentAccount.(*model.Account)
	if !ok {
		log.Fatalf("type assertion to *model.Account failed, got %T", ctx.CurrentAccount)
		return nil, err
	}

	return currentAccount, nil
}

// UpdateAccount updates an account
func (s *AccountService) UpdateAccount(ctx *core.Context, payload *types.UpdateAccountRequest) (m *model.Account, err error) {
	currentAccount, err := s.GetAccount(ctx)
	if err != nil {
		return nil, err
	}
	return s.repo.Update(currentAccount.ID, payload)
}
