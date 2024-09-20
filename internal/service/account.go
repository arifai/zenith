package service

import (
	"context"
	"errors"
	"github.com/arifai/zenith/config"
	"github.com/arifai/zenith/internal/model"
	"github.com/arifai/zenith/internal/repository"
	"github.com/arifai/zenith/internal/types/request"
	"github.com/arifai/zenith/internal/types/response"
	"github.com/arifai/zenith/pkg/crypto"
	"github.com/arifai/zenith/pkg/errormessage"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type (
	// AccountService provides methods to handle account-related operations in the application.
	AccountService interface {
		// Register handles the registration process for a new account by saving user data and hashed password in the database.
		Register(body *request.AccountCreateRequest) (*model.Account, error)

		// Authorization authenticates a user by validating their email and password, returning access and refresh tokens.
		Authorization(body *request.AccountAuthRequest) (*response.AccountAuthResponse, error)

		// Unauthorization invalidates both the access and refresh tokens present in the request body by blacklisting them.
		Unauthorization(body *request.AccountUnauthRequest) error

		// RefreshToken refreshes the access and refresh tokens for a given account ID if the provided refresh token is valid.
		RefreshToken(id uuid.UUID, body *request.AccountRefreshTokenRequest) (*response.AccountAuthResponse, error)

		// GetCurrent retrieves the current account details by the given account ID (uuid.UUID).
		GetCurrent(id *uuid.UUID) (*model.Account, error)

		// Update updates an existing account's details such as FullName and Email based on the provided id and request body.
		Update(id *uuid.UUID, body *request.AccountUpdateRequest) (*model.Account, error)

		UpdatePassword(id *uuid.UUID, body *request.AccountUpdatePasswordRequest) error
	}

	// accountService handles account-related operations and interacts with the account repository.
	accountService struct {
		*Service
		accountRepo repository.AccountRepository
	}
)

// NewAccountService initializes and returns an AccountService instance with the provided Service and AccountRepository.
func NewAccountService(service *Service, accountRepo repository.AccountRepository) AccountService {
	return &accountService{
		Service:     service,
		accountRepo: accountRepo,
	}
}

func (a *accountService) Register(body *request.AccountCreateRequest) (*model.Account, error) {
	tx := a.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
		if tx.Error != nil {
			tx.Rollback()
		}
	}()

	accountExist, err := a.accountRepo.FindByEmail(body.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return nil, err
	}
	if accountExist != nil {
		tx.Rollback()
		return nil, errormessage.ErrEmailAlreadyExists
	}

	passwordHash, err := generatePasswordHash(body.Password, a.config.PasswordSalt)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	account := &model.Account{FullName: body.FullName, Email: body.Email}
	account.AccountPassHashed = &model.AccountPassHashed{AccountId: account.ID, PassHashed: passwordHash}

	if err := a.accountRepo.Create(account); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return account, nil
}

func (a *accountService) Authorization(body *request.AccountAuthRequest) (*response.AccountAuthResponse, error) {
	account, err := a.accountRepo.FindByEmail(body.Email)
	if err != nil {
		return nil, errormessage.ErrEmailAddressNotFound
	} else if account == nil {
		return nil, errormessage.ErrEmailAddressNotFound
	}

	if err := validateAccount(account, body.Password); err != nil {
		return nil, err
	}

	accessToken, err := generateToken(account.ID, crypto.AccessToken, time.Hour*24)
	if err != nil {
		return nil, err
	}

	refreshToken, err := generateToken(account.ID, crypto.RefreshToken, time.Hour*24*30)
	if err != nil {
		return nil, err
	}

	return &response.AccountAuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (a *accountService) Unauthorization(body *request.AccountUnauthRequest) error {
	verifyAccessToken, err := crypto.VerifyToken(body.AccessToken, config.PublicKey)
	if err != nil {
		return errormessage.ErrInvalidAccessTokenInBody
	}

	if err = a.blacklistToken(verifyAccessToken.Jti.String(), verifyAccessToken.ExpiresAt); err != nil {
		return err
	}

	verifyRefreshToken, err := crypto.VerifyToken(body.RefreshToken, config.PublicKey)
	if err != nil {
		return errormessage.ErrInvalidRefreshTokenInBody
	}

	if err = a.blacklistToken(verifyRefreshToken.Jti.String(), verifyRefreshToken.ExpiresAt); err != nil {
		return err
	}

	return nil
}

func (a *accountService) RefreshToken(id uuid.UUID, body *request.AccountRefreshTokenRequest) (*response.AccountAuthResponse, error) {
	accessToken, err := generateToken(id, crypto.AccessToken, time.Hour*24)
	if err != nil {
		return nil, err
	}

	refreshToken, err := generateToken(id, crypto.RefreshToken, time.Hour*24*30)
	if err != nil {
		return nil, err
	}

	verifyRefreshToken, err := crypto.VerifyToken(body.RefreshToken, config.PublicKey)
	if err != nil {
		return nil, errormessage.ErrInvalidRefreshTokenInBody
	}

	if err = a.blacklistToken(verifyRefreshToken.Jti.String(), verifyRefreshToken.ExpiresAt); err != nil {
		return nil, err
	}

	return &response.AccountAuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (a *accountService) GetCurrent(id *uuid.UUID) (*model.Account, error) {
	account, err := a.accountRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (a *accountService) Update(id *uuid.UUID, body *request.AccountUpdateRequest) (*model.Account, error) {
	account, err := a.accountRepo.FindByID(id)
	if err != nil {
		return nil, errormessage.ErrAccountNotFound
	}

	account.FullName = body.FullName
	account.Email = body.Email

	if err := a.accountRepo.Update(account); err != nil {
		return nil, err
	}

	return account, nil
}

func (a *accountService) UpdatePassword(id *uuid.UUID, body *request.AccountUpdatePasswordRequest) error {
	account, err := a.accountRepo.FindByID(id)
	if err != nil {
		return errormessage.ErrAccountNotFound
	}

	if err := validateAccount(account, body.OldPassword); err != nil {
		return errormessage.ErrWrongOldPassword
	}

	passwordHash, err := generatePasswordHash(body.NewPassword, a.config.PasswordSalt)
	if err != nil {
		return err
	}

	accountPassHashed := &model.AccountPassHashed{AccountId: account.ID, PassHashed: passwordHash}
	account.AccountPassHashed = accountPassHashed

	if err := a.accountRepo.UpdatePassword(account); err != nil {
		return err
	}

	return nil
}

// generatePasswordHash generates a secure password hash using the Argon2ID algorithm with the provided password and salt.
// Returns the password hash as a string or an error if the hashing process fails.
func generatePasswordHash(password, salt string) (string, error) {
	return crypto.DefaultArgon2IDHash.GenerateHash([]byte(password), []byte(salt))
}

// validateAccount validates an account by checking if it's active and verifies the provided password. Returns an error on failure.
func validateAccount(account *model.Account, password string) error {
	if !account.Active {
		return errormessage.ErrAccountNotActive
	} else if account.AccountPassHashed == nil {
		return errormessage.ErrAccountPasswordHashMissing
	}

	valid, err := crypto.VerifyHash(password, account.AccountPassHashed.PassHashed)
	if err != nil {
		return err
	} else if !valid {
		return errormessage.ErrIncorrectPassword
	}

	return nil
}

// generateToken creates a signed token for a given accountID, tokenType, and duration.
// The generated token is returned as a string.
// In case of failure to generate an access or refresh token, an appropriate error is returned.
func generateToken(accountID uuid.UUID, tokenType string, duration time.Duration) (string, error) {
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
		case crypto.AccessToken:
			return "", errormessage.ErrFailedToGenerateAccessToken
		case crypto.RefreshToken:
			return "", errormessage.ErrFailedToGenerateRefreshToken
		}
	}

	return token, nil
}

// blacklistToken adds a token's "jti" to the blacklist with the specified expiration time in Redis.
func (a *accountService) blacklistToken(jti string, exp time.Time) error {
	ttl := time.Until(exp)

	if err := a.redis.Set(context.Background(), jti, "blacklisted", ttl).Err(); err != nil {
		return err
	}

	return nil
}
