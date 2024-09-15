package middleware

import (
	"errors"
	"github.com/arifai/go-modular-monolithic/config"
	"github.com/arifai/go-modular-monolithic/internal/account/domain/model"
	"github.com/arifai/go-modular-monolithic/internal/account/domain/repository"
	"github.com/arifai/go-modular-monolithic/pkg/common"
	crp "github.com/arifai/go-modular-monolithic/pkg/crypto"
	"github.com/arifai/go-modular-monolithic/pkg/errormessage"
	"github.com/arifai/go-modular-monolithic/pkg/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strings"
)

// Middleware is a Gin middleware for validating access tokens from Authorization headers and setting authorized account context.
func Middleware(db *gorm.DB) gin.HandlerFunc {
	repo := repository.NewAccountRepository(db)
	return func(ctx *gin.Context) {
		resp := common.Response{}
		if account, err := validateAndExtractAccount(ctx, repo); err != nil {
			resp.Unauthorized(ctx, []utils.IError{}, err.Error())
			ctx.Abort()
		} else {
			ctx.Set("account", account)
			ctx.Next()
		}
	}
}

// validateAndExtractAccount validates the token and extracts the account information.
func validateAndExtractAccount(ctx *gin.Context, repo *repository.AccountRepository) (*model.Account, error) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		return nil, errors.New(errormessage.ErrMissingAuthorizationHeaderText)
	}

	tokenString, err := extractToken(authHeader)
	if err != nil {
		return nil, err
	}

	tokenPayload, err := crp.VerifyToken(tokenString, config.PublicKey)
	if err != nil {
		return nil, err
	}

	if tokenPayload.TokenType != "access_token" {
		return nil, errors.New(errormessage.ErrInvalidTokenTypeText)
	}

	account, err := repo.Find(tokenPayload.AccountId)
	if err != nil {
		return nil, errors.New(errormessage.ErrCannotFindAuthorizedAccountText)
	}
	return account, nil
}

// extractToken splits the authorization header to extract the token.
func extractToken(authHeader string) (string, error) {
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return "", errors.New(errormessage.ErrInvalidAccessTokenText)
	}
	return tokenParts[1], nil
}
