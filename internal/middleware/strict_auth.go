package middleware

import (
	"context"
	"errors"
	"github.com/arifai/zenith/config"
	"github.com/arifai/zenith/pkg/common"
	"github.com/arifai/zenith/pkg/crypto"
	"github.com/arifai/zenith/pkg/errormessage"
	"github.com/arifai/zenith/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"strings"
)

// StrictAuthMiddleware struct provides methods for strict authorization checks and token blacklisting using a Redis backend.
type StrictAuthMiddleware struct{ *Middleware }

func NewStrictAuthMiddleware(middleware *Middleware) *StrictAuthMiddleware {
	return &StrictAuthMiddleware{middleware}
}

// StrictAuth is a middleware function that validates and extracts the account from the authorization header.
func (s *StrictAuthMiddleware) StrictAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var response common.Response
		id, err := s.validateAndExtractAccount(ctx)
		if err != nil {
			response.Unauthorized(ctx, []utils.IError{}, err.Error())
			ctx.Abort()
			return
		} else {
			ctx.Set("account_id", id)
			ctx.Next()
			return
		}
	}
}

// IsTokenBlacklisted checks if a given token's jti is present in the Redis blacklist.
func (s *StrictAuthMiddleware) IsTokenBlacklisted(jti string) (bool, error) {
	value, err := s.getRedisValue(jti)
	if errors.Is(err, redis.Nil) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return value == "blacklisted", nil
}

// validateAndExtractAccount validates the authorization header and extracts the associated account.
func (s *StrictAuthMiddleware) validateAndExtractAccount(ctx *gin.Context) (*uuid.UUID, error) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		return nil, errormessage.ErrMissingAuthorizationHeader
	}

	tokenString, err := s.extractToken(authHeader)
	if err != nil {
		return nil, err
	}

	tokenPayload, err := crypto.VerifyToken(tokenString, config.PublicKey)
	if err != nil {
		return nil, err
	} else if tokenPayload.TokenType != crypto.AccessToken {
		return nil, errormessage.ErrInvalidTokenType
	}

	isTokenBlacklisted, err := s.IsTokenBlacklisted(tokenPayload.Jti.String())
	if err != nil {
		return nil, err
	} else if isTokenBlacklisted {
		return nil, errormessage.ErrInvalidAccessToken
	}

	return &tokenPayload.AccountID, nil
}

// extractToken splits the authorization header to extract the token.
func (s *StrictAuthMiddleware) extractToken(authHeader string) (string, error) {
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return "", errormessage.ErrInvalidAccessToken
	}

	return tokenParts[1], nil
}

func (s *StrictAuthMiddleware) getRedisValue(jti string) (string, error) {
	result, err := s.redis.Get(context.Background(), jti).Result()
	if err != nil {
		return "", err
	}

	return result, nil
}
