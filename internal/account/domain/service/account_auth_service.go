package service

import (
	"github.com/arifai/go-modular-monolithic/config"
	"github.com/arifai/go-modular-monolithic/internal/account/api/types"
	"github.com/arifai/go-modular-monolithic/internal/account/domain/repository"
	"github.com/arifai/go-modular-monolithic/pkg/common"
	"github.com/arifai/go-modular-monolithic/pkg/crypto"
	"github.com/arifai/go-modular-monolithic/pkg/errormessage"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// AccountAuthService is used for handling account authorization processes.
type AccountAuthService struct {
	config *config.Config
	repo   *repository.AccountRepository
}

// NewAccountAuthService initializes an AccountAuthService with a given database connection and configuration settings.
func NewAccountAuthService(db *gorm.DB, config *config.Config) *AccountAuthService {
	return &AccountAuthService{
		config: config,
		repo:   repository.NewAccountRepository(db),
	}
}

// Authorize authenticates an account using the provided email and password.
// It returns a common.AuthResponse containing access and refresh tokens or an error.
func (a AccountAuthService) Authorize(payload *types.AccountAuthRequest) (*common.AuthResponse, error) {
	account, err := a.repo.FindByEmail(payload.Email)
	if err != nil {
		return nil, errormessage.ErrEmailAddressNotFound
	} else if !account.Active {
		return nil, errormessage.ErrAccountNotActive
	} else if account.AccountPassHashed == nil {
		return nil, errormessage.ErrAccountPasswordHashMissing
	}

	valid, err := crypto.VerifyHash(payload.Password, account.AccountPassHashed.PassHashed)
	if err != nil {
		return nil, err
	} else if !valid {
		return nil, errormessage.ErrIncorrectPassword
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
		return nil, errormessage.ErrFailedToGenerateAccessToken
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
		return nil, errormessage.ErrFailedToGenerateRefreshToken
	}

	return &common.AuthResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
