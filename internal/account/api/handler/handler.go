package handler

import (
	"github.com/arifai/go-modular-monolithic/config"
	"github.com/arifai/go-modular-monolithic/internal/account/api/types"
	"github.com/arifai/go-modular-monolithic/internal/account/domain/service"
	"github.com/arifai/go-modular-monolithic/pkg/common"
	"github.com/arifai/go-modular-monolithic/pkg/core"
	"github.com/arifai/go-modular-monolithic/pkg/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuthHandler handles the authentication requests by validating input, using the account service to authorize,
// and responding accordingly.
func AuthHandler(ctx *gin.Context, db *gorm.DB, config *config.Config) {
	resp := new(common.Response)
	accountService := service.NewAccountAuthService(db, config)
	body, err := utils.ValidateBody[types.AccountAuthRequest](ctx)
	if err != nil {
		resp.Error(ctx, err)
		return
	}

	result, err := accountService.Authorize(body)
	if err != nil {
		resp.Error(ctx, err)
		return
	}

	resp.Authorized(ctx, result)
}

// GetAccountHandler handles the retrieval of an account using the GetAccount service.
// It initializes the account service and context, then attempts to get the account.
// Upon success, it sends a success response; otherwise, it sends an error response.
func GetAccountHandler(ctx *gin.Context, db *gorm.DB, config *config.Config) {
	resp := new(common.Response)
	accountService := service.NewAccountService(db, config)
	context := core.NewContext(ctx)
	result, err := accountService.GetAccount(context)
	if err != nil {
		resp.Error(ctx, err)
		return
	}

	resp.Success(ctx, result)
}

// RegisterAccountHandler handles the registration of a new user account.
// It validates the request body, creates a new user account using the CreateAccount service,
// and returns an appropriate response based on the operation outcome.
func RegisterAccountHandler(ctx *gin.Context, db *gorm.DB, config *config.Config) {
	resp := new(common.Response)
	accountService := service.NewAccountService(db, config)
	body, err := utils.ValidateBody[types.AccountCreateRequest](ctx)
	if err != nil {
		resp.Error(ctx, err)
		return
	}

	result, err := accountService.CreateAccount(body)
	if err != nil {
		resp.Error(ctx, err)
		return
	}

	resp.Created(ctx, "", result)
}

// UpdateAccountHandler handles the HTTP request for updating an account.
// It validates the request body, invokes the account service to update the account,
// and sends the appropriate response based on the success or failure of the update operation.
func UpdateAccountHandler(ctx *gin.Context, db *gorm.DB, config *config.Config) {
	resp := new(common.Response)
	accountService := service.NewAccountService(db, config)
	context := core.NewContext(ctx)
	body, err := utils.ValidateBody[types.AccountUpdateRequest](ctx)
	if err != nil {
		resp.Error(ctx, err)
		return
	}

	result, err := accountService.UpdateAccount(context, body)
	if err != nil {
		resp.Error(ctx, err)
		return
	}

	resp.Success(ctx, result)
}
