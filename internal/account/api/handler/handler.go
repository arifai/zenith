package handler

import (
	"github.com/arifai/zenith/config"
	"github.com/arifai/zenith/internal/account/api/types"
	"github.com/arifai/zenith/internal/account/domain/service"
	"github.com/arifai/zenith/pkg/common"
	"github.com/arifai/zenith/pkg/core"
	"github.com/arifai/zenith/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// AccountHandler manages account-related operations such as authentication, retrieval, registration, and updates.
// It uses a gorm.DB instance for database operations, config.Config for configuration settings, and common.Response for responses.
type AccountHandler struct {
	db          *gorm.DB
	config      *config.Config
	redisClient *redis.Client
	resp        *common.Response
}

// NewAccountHandler initializes an AccountHandler with a given database connection and configuration settings.
func NewAccountHandler(db *gorm.DB, config *config.Config, redisClient *redis.Client) *AccountHandler {
	return &AccountHandler{
		db:          db,
		config:      config,
		redisClient: redisClient,
		resp:        new(common.Response),
	}
}

// AuthHandler handles the authentication logic for user accounts, including request validation and response formatting.
func (a *AccountHandler) AuthHandler(ctx *gin.Context) {
	accountAuthService := service.NewAccountAuthService(a.db, a.config, a.redisClient)
	body, err := utils.ValidateBody[types.AccountAuthRequest](ctx)
	if err != nil {
		a.resp.Error(ctx, err)
		return
	}

	result, err := accountAuthService.Authorize(body)
	if err != nil {
		a.resp.Error(ctx, err)
		return
	}

	a.resp.Authorized(ctx, result)
}

// UnauthHandler handles the unauthentication of user tokens by validating the request body and processing the tokens.
func (a *AccountHandler) UnauthHandler(ctx *gin.Context) {
	accountAuthService := a.initializeAccountAuthService()

	body, err := utils.ValidateBody[types.AccountUnauthRequest](ctx)
	if err != nil {
		a.resp.Error(ctx, err)
		return
	}

	if err := a.handleUnauthorizedTokens(accountAuthService, body); err != nil {
		a.resp.Error(ctx, err)
		return
	}

	a.resp.Success(ctx, nil)
}

// initializeAccountAuthService initializes and returns an AccountAuthService instance using the AccountHandler's dependencies.
func (a *AccountHandler) initializeAccountAuthService() *service.AccountAuthService {
	return service.NewAccountAuthService(a.db, a.config, a.redisClient)
}

// handleUnauthorizedTokens invalidates both access and refresh tokens using the AccountAuthService and returns an error if any.
func (a *AccountHandler) handleUnauthorizedTokens(accountAuthService *service.AccountAuthService, request *types.AccountUnauthRequest) error {
	if err := accountAuthService.Unauthorized(request.AccessToken); err != nil {
		return err
	}

	if err := accountAuthService.Unauthorized(request.RefreshToken); err != nil {
		return err
	}

	return nil
}

// GetAccountHandler retrieves the current user's account details. It initializes account service and context,
// fetches account information, and sends a success or error response based on the result.
func (a *AccountHandler) GetAccountHandler(ctx *gin.Context) {
	accountService := service.NewAccountService(a.db, a.config, a.redisClient)
	context := core.NewContext(ctx)
	result, err := accountService.GetAccount(context)
	if err != nil {
		a.resp.Error(ctx, err)
		return
	}

	a.resp.Success(ctx, result)
}

// RegisterAccountHandler handles the registration of a new user account by validating the request
// body, invoking the account service to create the account, and sending an appropriate response.
func (a *AccountHandler) RegisterAccountHandler(ctx *gin.Context) {
	accountService := service.NewAccountService(a.db, a.config, a.redisClient)
	body, err := utils.ValidateBody[types.AccountCreateRequest](ctx)
	if err != nil {
		a.resp.Error(ctx, err)
		return
	}

	result, err := accountService.CreateAccount(body, a.config)
	if err != nil {
		a.resp.Error(ctx, err)
		return
	}

	a.resp.Created(ctx, "", result)
}

// UpdateAccountHandler handles the updating of user account information. It validates the request body, updates
// the account via the account service, and sends a success or error response based on the result.
func (a *AccountHandler) UpdateAccountHandler(ctx *gin.Context) {
	accountService := service.NewAccountService(a.db, a.config, a.redisClient)
	context := core.NewContext(ctx)
	body, err := utils.ValidateBody[types.AccountUpdateRequest](ctx)
	if err != nil {
		a.resp.Error(ctx, err)
		return
	}

	result, err := accountService.UpdateAccount(context, body)
	if err != nil {
		a.resp.Error(ctx, err)
		return
	}

	a.resp.Success(ctx, result)
}

// UpdatePasswordAccountHandler handles the password update request for an account.
// It validates the request body, calls the account service to update the password, and sends the appropriate response.
func (a *AccountHandler) UpdatePasswordAccountHandler(ctx *gin.Context) {
	accountService := service.NewAccountService(a.db, a.config, a.redisClient)
	context := core.NewContext(ctx)
	body, err := utils.ValidateBody[types.AccountUpdatePasswordRequest](ctx)
	if err != nil {
		a.resp.Error(ctx, err)
		return
	}

	result, err := accountService.UpdatePassword(context, body, a.config)
	if err != nil {
		a.resp.Error(ctx, err)
		return
	}

	a.resp.Success(ctx, result)
}
