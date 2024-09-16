package service

import (
	"context"
	"github.com/arifai/zenith/config"
	"github.com/arifai/zenith/internal/account/api/types"
	"github.com/arifai/zenith/internal/account/domain/repository"
	"github.com/arifai/zenith/pkg/common"
	"github.com/arifai/zenith/pkg/crypto"
	"github.com/arifai/zenith/pkg/errormessage"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"time"
)

// AccountAuthService is used for handling account authorization processes.
type AccountAuthService struct {
	defaultConfig *config.Config
	redisClient   *redis.Client
	repo          *repository.AccountRepository
}

// NewAccountAuthService initializes an AccountAuthService with a given database connection and configuration settings.
func NewAccountAuthService(db *gorm.DB, defaultConfig *config.Config, redisClient *redis.Client) *AccountAuthService {
	return &AccountAuthService{
		defaultConfig: defaultConfig,
		redisClient:   redisClient,
		repo:          repository.NewAccountRepository(db, redisClient),
	}
}

// Authorize authenticates an account using the provided email and password.
// It returns a common.AuthResponse containing access and refresh tokens or an error.
func (a AccountAuthService) Authorize(payload *types.AccountAuthRequest) (*common.AuthResponse, error) {
	account, err := a.repo.FindByEmail(payload.Email)
	if err != nil {
		return nil, errormessage.ErrEmailAddressNotFound
	}
	if !account.Active {
		return nil, errormessage.ErrAccountNotActive
	}
	if account.AccountPassHashed == nil {
		return nil, errormessage.ErrAccountPasswordHashMissing
	}
	valid, err := crypto.VerifyHash(payload.Password, account.AccountPassHashed.PassHashed)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, errormessage.ErrIncorrectPassword
	}

	accessToken, err := a.generateToken(account.ID, "access_token", time.Hour*24)
	if err != nil {
		return nil, err
	}

	refreshToken, err := a.generateToken(account.ID, "refresh_token", time.Hour*168)
	if err != nil {
		return nil, err
	}

	return &common.AuthResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

// Unauthorized invalidates an access token by verifying it and adding it to the blacklist.
func (a AccountAuthService) Unauthorized(token string) error {
	verifyToken, err := crypto.VerifyToken(token, config.PublicKey)
	if err != nil {
		return err
	}

	if err = a.blacklistToken(verifyToken.Jti.String(), verifyToken.ExpiresAt); err != nil {
		return err
	}

	return nil
}

// blacklistToken blacklists a Paseto by storing its jti in Redis with an expiration time.
func (a AccountAuthService) blacklistToken(jti string, exp time.Time) error {
	ttl := time.Until(exp)
	err := a.redisClient.Set(context.Background(), jti, "blacklisted", ttl).Err()
	if err != nil {
		return err
	}

	return nil
}

// generateToken creates a token with specified type and duration
func (a AccountAuthService) generateToken(accountID uuid.UUID, tokenType string, duration time.Duration) (string, error) {
	now := time.Now()
	payload := crypto.TokenPayload{
		Jti:       uuid.New(),
		AccountId: accountID,
		IssuedAt:  now,
		NotBefore: now,
		ExpiresAt: now.Add(duration),
		TokenType: tokenType,
	}
	token := payload.GenerateToken(config.SecretKey)
	if token == "" {
		switch tokenType {
		case "access_token":
			return "", errormessage.ErrFailedToGenerateAccessToken
		case "refresh_token":
			return "", errormessage.ErrFailedToGenerateRefreshToken
		}
	}
	return token, nil
}
