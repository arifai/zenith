package service

import (
	"errors"
	"github.com/arifai/zenith/config"
	"github.com/arifai/zenith/internal/model"
	"github.com/arifai/zenith/internal/repository"
	"github.com/arifai/zenith/internal/types/request"
	"github.com/arifai/zenith/internal/types/response"
	"github.com/arifai/zenith/pkg/crypto"
	"github.com/arifai/zenith/pkg/errormessage"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strings"
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
		RefreshToken(body *request.AccountRefreshTokenRequest) (*response.AccountAuthResponse, error)

		// GetCurrent retrieves the current account details by the given account ID (uuid.UUID).
		GetCurrent(id *uuid.UUID) (*model.Account, error)

		// Update updates an existing account's details such as FullName and Email based on the provided id and request body.
		Update(id *uuid.UUID, body *request.AccountUpdateRequest) (*model.Account, error)

		// UpdatePassword updates the password of an account identified by the given UUID. It takes the new password and the old password for validation. Returns an error if the operation fails.
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
	return &accountService{Service: service, accountRepo: accountRepo}
}

func (a *accountService) Register(body *request.AccountCreateRequest) (*model.Account, error) {
	founded, err := a.accountRepo.FindByEmail(body.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if founded != nil {
		return nil, errormessage.ErrEmailAlreadyExists
	}

	newAccount := &model.Account{FullName: body.FullName, Email: strings.ToLower(body.Email)}
	passwordHash, err := generatePasswordHash(body.Password, a.PasswordSalt)
	if err != nil {
		return nil, err
	}

	if err := a.accountRepo.Create(newAccount, passwordHash); err != nil {
		return nil, err
	}

	return newAccount, nil
}

func (a *accountService) Authorization(body *request.AccountAuthRequest) (*response.AccountAuthResponse, error) {
	account, err := a.accountRepo.FindByEmail(body.Email)
	if err != nil {
		return nil, errormessage.ErrEmailAddressNotFound
	}

	if err = a.accountRepo.SetFCMToken(account.Email, body.FcmToken); err != nil {
		return nil, err
	}

	if err := validateAccount(account, body.Password); err != nil {
		return nil, err
	}

	parsedDeviceID, err := uuid.Parse(body.DeviceID)
	if err != nil {
		a.Logger.Error(errormessage.ErrFailedToParseUUIDText, zap.String("input", body.DeviceID), zap.Error(err))
		return nil, errormessage.ErrInvalidDeviceIDInBody
	}

	accessToken, err := generateToken(account.ID, parsedDeviceID, crypto.AccessToken, time.Hour*6)
	if err != nil {
		return nil, err
	}

	refreshToken, err := generateToken(account.ID, parsedDeviceID, crypto.RefreshToken, time.Hour*24*30)
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

	verifyRefreshToken, err := crypto.VerifyToken(body.RefreshToken, config.PublicKey)
	if err != nil {
		return errormessage.ErrInvalidRefreshTokenInBody
	}

	if err = a.accountRepo.BlacklistToken(verifyAccessToken.Jti.String(), verifyAccessToken.ExpiresAt); err != nil {
		return err
	}

	if err = a.accountRepo.BlacklistToken(verifyRefreshToken.Jti.String(), verifyRefreshToken.ExpiresAt); err != nil {
		return err
	}

	if err = a.accountRepo.UnsetFCMToken(verifyAccessToken.AccountId); err != nil {
		return err
	}

	return nil
}

func (a *accountService) RefreshToken(body *request.AccountRefreshTokenRequest) (*response.AccountAuthResponse, error) {
	verifyRefreshToken, err := crypto.VerifyToken(body.RefreshToken, config.PublicKey)
	if err != nil {
		return nil, errormessage.ErrInvalidRefreshTokenInBody
	}

	if err = a.accountRepo.BlacklistToken(verifyRefreshToken.Jti.String(), verifyRefreshToken.ExpiresAt); err != nil {
		return nil, err
	}

	accessToken, err := generateToken(verifyRefreshToken.AccountId, verifyRefreshToken.DeviceID, crypto.AccessToken, time.Hour*6)
	if err != nil {
		return nil, err
	}

	refreshToken, err := generateToken(verifyRefreshToken.AccountId, verifyRefreshToken.DeviceID, crypto.RefreshToken, time.Hour*24*30)
	if err != nil {
		return nil, err
	}

	return &response.AccountAuthResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
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
	account.Email = strings.ToLower(body.Email)

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
		return err
	}

	passwordHash, err := generatePasswordHash(body.NewPassword, a.PasswordSalt)
	if err != nil {
		return err
	}

	account.AccountPassHashed = &model.AccountPassHashed{AccountID: account.ID, PassHashed: passwordHash}

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
	}

	valid, err := crypto.VerifyHash(password, account.AccountPassHashed.PassHashed)
	if err != nil {
		return err
	} else if !valid {
		return errormessage.ErrWrongOldPassword
	}

	return nil
}

// generateToken creates a signed token for a given accountID, deviceID, tokenType, and duration.
// The generated token is returned as a string.
// In case of failure to generate an access or refresh token, an appropriate error is returned.
func generateToken(accountID, deviceID uuid.UUID, tokenType string, duration time.Duration) (string, error) {
	now := time.Now()
	payload := crypto.TokenPayload{
		Jti:       uuid.New(),
		DeviceID:  deviceID,
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
