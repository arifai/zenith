package handler

import (
	"github.com/arifai/zenith/internal/service"
	"github.com/arifai/zenith/internal/types/request"
	"github.com/arifai/zenith/pkg/utils"
	"github.com/gin-gonic/gin"
)

// AccountHandler handles HTTP requests related to account operations such as registration, authorization, and updates.
type AccountHandler struct {
	*Handler
	accountService service.AccountService
}

// NewAccountHandler initializes a new AccountHandler with the provided Handler and AccountService.
func NewAccountHandler(handler *Handler, accountService service.AccountService) *AccountHandler {
	return &AccountHandler{Handler: handler, accountService: accountService}
}

// Register handles HTTP requests for creating a new user account.
// It validates the request body, calls the account service to register the account,
// and sends appropriate HTTP responses based on the outcome.
func (a *AccountHandler) Register(ctx *gin.Context) {
	body, err := utils.ValidateBody[request.AccountCreateRequest](ctx)
	if err != nil {
		a.response.Error(ctx, err)
		return
	}

	result, err := a.accountService.Register(body)
	if err != nil {
		a.response.Error(ctx, err)
		return
	}

	a.response.Created(ctx, "Account successfully registered", result)
}

// Authorization handles the authorization of an account request by validating the request body,
// invoking the account service's Authorization method, and sending the appropriate HTTP responses.
func (a *AccountHandler) Authorization(ctx *gin.Context) {
	body, err := utils.ValidateBody[request.AccountAuthRequest](ctx)
	if err != nil {
		a.response.Error(ctx, err)
		return
	}

	result, err := a.accountService.Authorization(body)
	if err != nil {
		a.response.Error(ctx, err)
		return
	}

	a.response.Authorized(ctx, result)
}

// Unauthorization provides HTTP handling for unauthorizing an account by invalidating access and refresh tokens.
func (a *AccountHandler) Unauthorization(ctx *gin.Context) {
	body, err := utils.ValidateBody[request.AccountUnauthRequest](ctx)
	if err != nil {
		a.response.Error(ctx, err)
		return
	}

	if err := a.accountService.Unauthorization(body); err != nil {
		a.response.Error(ctx, err)
		return
	}

	a.response.Success(ctx, nil)
}

// RefreshToken handles the HTTP request to refresh an authentication token.
// It validates the request body and invokes the account service's RefreshToken method.
func (a *AccountHandler) RefreshToken(ctx *gin.Context) {
	body, err := utils.ValidateBody[request.AccountRefreshTokenRequest](ctx)
	if err != nil {
		a.response.Error(ctx, err)
		return
	}

	result, err := a.accountService.RefreshToken(body)
	if err != nil {
		a.response.Error(ctx, err)
		return
	}

	a.response.Success(ctx, result)
}

// GetCurrent handles the retrieval of the current account details based on the account ID from the context.
func (a *AccountHandler) GetCurrent(ctx *gin.Context) {
	accountId := GetAccountIDFromContext(ctx)
	if accountId == nil {
		a.response.NotFound(ctx, "Account ID not found in context")
		return
	}

	result, err := a.accountService.GetCurrent(accountId)
	if err != nil {
		a.response.Error(ctx, err)
		return
	}

	a.response.Success(ctx, result)
}

// Update handles the updating of an account's information.
// It validates the request body, calls the account service to update the account,
// and sends appropriate HTTP responses based on the outcome.
func (a *AccountHandler) Update(ctx *gin.Context) {
	accountId := GetAccountIDFromContext(ctx)
	body, err := utils.ValidateBody[request.AccountUpdateRequest](ctx)
	if err != nil {
		a.response.Error(ctx, err)
		return
	}

	account, err := a.accountService.Update(accountId, body)
	if err != nil {
		a.response.Error(ctx, err)
		return
	}

	a.response.Success(ctx, account)
}

// UpdatePassword handles the HTTP request to update an account's password.
// It retrieves the account ID from the context and validates the request body.
// If validation passes, it calls the account service to update the password and sends an appropriate response.
func (a *AccountHandler) UpdatePassword(ctx *gin.Context) {
	accountId := GetAccountIDFromContext(ctx)
	body, err := utils.ValidateBody[request.AccountUpdatePasswordRequest](ctx)
	if err != nil {
		a.response.Error(ctx, err)
		return
	}

	if err := a.accountService.UpdatePassword(accountId, body); err != nil {
		a.response.Error(ctx, err)
		return
	}

	a.response.Success(ctx, nil)
}
