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

// AccountService provides methods to manage user accounts
// It integrates configuration settings and account repository for CRUD operations
type AccountService struct {
	config *config.Config
	repo   *repository.AccountRepository
}

// NewAccountService initializes a new AccountService with the given database context and configuration settings.
func NewAccountService(db *gorm.DB, config *config.Config) *AccountService {
	return &AccountService{
		config: config,
		repo:   repository.NewAccountRepository(db),
	}
}

// CreateAccount registers a new user account in the system using the provided payload data.
// The payload must contain full name, email, and password. Returns the created model.Account or any error encountered.
func (s *AccountService) CreateAccount(payload *types.AccountCreateRequest) (*model.Account, error) {
	return s.repo.CreateAccount(payload)
}

// GetAccount retrieves the current account from the given context.
// It casts ctx.CurrentAccount to a model.Account pointer.
// Returns the current model.Account or an error if the type assertion fails.
func (s *AccountService) GetAccount(ctx *core.Context) (m *model.Account, err error) {
	currentAccount, ok := ctx.CurrentAccount.(*model.Account)
	if !ok {
		log.Fatalf("type assertion to *model.Account failed, got %T", ctx.CurrentAccount)
		return nil, err
	}

	return currentAccount, nil
}

// UpdateAccount updates the details of the current account using the provided payload.
// It first retrieves the current account from the context, then updates the account data in the repository.
// Returns the updated model.Account or an error if encountered.
func (s *AccountService) UpdateAccount(ctx *core.Context, payload *types.AccountUpdateRequest) (m *model.Account, err error) {
	currentAccount, err := s.GetAccount(ctx)
	if err != nil {
		return nil, err
	}
	return s.repo.Update(currentAccount.ID, payload)
}
