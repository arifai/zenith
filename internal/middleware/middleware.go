package middleware

import (
	"github.com/arifai/go-modular-monolithic/config"
	"github.com/arifai/go-modular-monolithic/internal/account/domain/repository"
	errmsg "github.com/arifai/go-modular-monolithic/internal/errors"
	"github.com/arifai/go-modular-monolithic/pkg/common"
	crp "github.com/arifai/go-modular-monolithic/pkg/crypto"
	"github.com/arifai/go-modular-monolithic/pkg/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strings"
)

// Middleware is a Gin middleware for validating access tokens from Authorization headers and setting authorized account context.
func Middleware(db *gorm.DB) gin.HandlerFunc {
	repo := repository.NewAccountRepository(db)

	return func(c *gin.Context) {
		resp := common.Response{}
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			resp.Unauthorized(c, []utils.IError{}, errmsg.ErrMissingAuthorizationHeaderText)
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			resp.Unauthorized(c, []utils.IError{}, errmsg.ErrInvalidAccessTokenText)
			c.Abort()
			return
		}

		tokenString := tokenParts[1]
		tokenPayload, err := crp.VerifyToken(tokenString, config.PublicKey)
		if err != nil {
			resp.Unauthorized(c, []utils.IError{}, err.Error())
			c.Abort()
			return
		}
		if tokenPayload.TokenType != "access_token" {
			resp.Unauthorized(c, []utils.IError{}, errmsg.ErrInvalidTokenTypeText)
			c.Abort()
			return
		}

		account, err := repo.Find(tokenPayload.AccountId)
		if err != nil {
			resp.Unauthorized(c, []utils.IError{}, errmsg.ErrCannotFindAuthorizedAccountText)
			c.Abort()
			return
		}

		c.Set("account", account)
		c.Next()
	}
}
