package service

import (
	"github.com/arifai/go-modular-monolithic/config"
	"github.com/arifai/go-modular-monolithic/internal/account/api/types"
	"github.com/arifai/go-modular-monolithic/internal/account/domain/repository"
	errmsg "github.com/arifai/go-modular-monolithic/internal/errors"
	"github.com/arifai/go-modular-monolithic/pkg/common"
	"github.com/arifai/go-modular-monolithic/pkg/crypto"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// AccountAuthService is a struct that represent the account auth repository
type AccountAuthService struct {
	config *config.Config
	repo   *repository.AccountRepository
}

// NewAccountAuthService creates a new account auth repository
func NewAccountAuthService(db *gorm.DB, config *config.Config) *AccountAuthService {
	return &AccountAuthService{
		config: config,
		repo:   repository.NewAccountRepository(db),
	}
}

// Authorize authorizes an account
func (a AccountAuthService) Authorize(body *types.AccountAuthRequest) (*common.AuthResponse, error) {
	account, err := a.repo.FindByEmail(body.Email)
	if err != nil {
		return nil, errmsg.ErrEmailAddressNotFound
	} else if !account.Active {
		return nil, errmsg.ErrAccountNotActive
	} else if account.AccountPassHashed == nil {
		return nil, errmsg.ErrAccountPasswordHashMissing
	}

	valid, err := crypto.VerifyHash(body.Password, account.AccountPassHashed.PassHashed)
	if err != nil {
		return nil, err
	} else if !valid {
		return nil, errmsg.ErrIncorrectPassword
	}

	tn := time.Now()
	accessTokenPayload := crypto.TokenPayload{
		Jti:       uuid.New(),
		AccountId: account.ID,
		IssuedAt:  tn,
		NotBefore: tn,
		ExpiresAt: tn.Add(time.Hour * 24),
		TokenType: "access_token",
	}

	accessToken := accessTokenPayload.GenerateToken(config.SecretKey)
	if accessToken == "" {
		return nil, errmsg.ErrFailedToGenerateAccessToken
	}

	refreshTokenPayload := crypto.TokenPayload{
		Jti:       uuid.New(),
		AccountId: account.ID,
		IssuedAt:  tn,
		NotBefore: tn,
		ExpiresAt: tn.Add(time.Hour * 168),
		TokenType: "refresh_token",
	}

	refreshToken := refreshTokenPayload.GenerateToken(config.SecretKey)
	if refreshToken == "" {
		return nil, errmsg.ErrFailedToGenerateRefreshToken
	}

	return &common.AuthResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
