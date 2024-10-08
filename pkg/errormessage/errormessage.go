package errormessage

import "errors"

const (
	ErrEmailAlreadyExistsText           = "email already exists"
	ErrEmailAddressNotFoundText         = "email address not found"
	ErrAccountNotActiveText             = "your account does not active"
	ErrAccountPasswordHashMissingText   = "account password hash is missing"
	ErrIncorrectPasswordText            = "incorrect password"
	ErrFailedToGenerateAccessTokenText  = "failed to generate access token"
	ErrFailedToGenerateRefreshTokenText = "failed to generate refresh token"
	ErrInvalidEncodedHashText           = "invalid encoded hash"
	ErrIncompatibleArgon2VersionText    = "incompatible argon2 version"
	ErrMissingAuthorizationHeaderText   = "authorization header missing"
	ErrFailedToConnectRedisText         = "failed to connect to Redis"
	ErrInvalidTokenHashText             = "invalid token hash"
	ErrParsingRequestDataText           = "failed to parsing request data"
	ErrFailedToConnectDBText            = "failed to connect to database"
	ErrFailedGetDBInstanceText          = "failed to get database instance"
	ErrFailedParseTokenText             = "failed to parse token"
	ErrFailedParsePublicHexText         = "failed to parse public hex"
	ErrFailedGetJTIText                 = "failed to to get 'jti'"
	ErrFailedParseJTIText               = "failed to parse 'jti'"
	ErrFailedGetSubText                 = "failed to get 'sub'"
	ErrFailedParseACIText               = "failed to parse account ID"
	ErrFailedGetIATText                 = "failed to get 'iat'"
	ErrFailedGetNBFText                 = "failed to get 'nbf'"
	ErrFailedGetEXPText                 = "failed to get 'exp'"
	ErrFailedGetAudText                 = "failed to get 'aud'"
	ErrFailedParseAudText               = "failed to parse 'aud'"
	ErrTokenExpiredText                 = "token has expired"
	ErrBadRequestText                   = "your request is invalid"
	ErrRequestBodyEmptyText             = "request body is empty"
	ErrInvalidTokenTypeText             = "invalid token type"
	ErrWrongOldPasswordText             = "wrong old password"
	ErrInvalidSaltLengthText            = "invalid salt length"
	ErrInvalidAccessTokenInBodyText     = "invalid access token in body, maybe your token has expired"
	ErrInvalidRefreshTokenInBodyText    = "invalid refresh token in body, maybe your token has expired"
	ErrAccountNotFoundText              = "account not found"
	ErrFailedSendEmailText              = "failed to send email"
	ErrInitializingServerText           = "error initializing server"
	ErrFailedSetTrustedProxiesText      = "failed to set trusted proxies"
	ErrSaltLengthIncorrectText          = "salt length is incorrect"
	ErrFailedToDecodeBase64Text         = "failed to decode base64"
	ErrFailedToGenerateRandomBytesText  = "failed to generate random bytes"
	ErrFailedToLoadEnvFileText          = "failed to load env file"
	ErrFailedToLoadEnvVariableText      = "failed to load env variable"
	ErrMigrationText                    = "error during migration"
	ErrCreatingEnumsText                = "error creating enums"
	ErrInsertingMigrationDataText       = "error during inserting migration data"
	ErrFailedToParseUUIDText            = "failed to parse uuid"
	ErrInvalidDeviceIDInBodyText        = "invalid device ID"
)

var (
	ErrEmailAlreadyExists           = errors.New(ErrEmailAlreadyExistsText)
	ErrEmailAddressNotFound         = errors.New(ErrEmailAddressNotFoundText)
	ErrAccountNotActive             = errors.New(ErrAccountNotActiveText)
	ErrAccountPasswordHashMissing   = errors.New(ErrAccountPasswordHashMissingText)
	ErrIncorrectPassword            = errors.New(ErrIncorrectPasswordText)
	ErrFailedToGenerateAccessToken  = errors.New(ErrFailedToGenerateAccessTokenText)
	ErrFailedToGenerateRefreshToken = errors.New(ErrFailedToGenerateRefreshTokenText)
	ErrInvalidEncodedHash           = errors.New(ErrInvalidEncodedHashText)
	ErrIncompatibleArgon2Version    = errors.New(ErrIncompatibleArgon2VersionText)
	ErrWrongOldPassword             = errors.New(ErrWrongOldPasswordText)
	ErrInvalidSaltLength            = errors.New(ErrInvalidSaltLengthText)
	ErrMissingAuthorizationHeader   = errors.New(ErrMissingAuthorizationHeaderText)
	ErrInvalidTokenType             = errors.New(ErrInvalidTokenTypeText)
	ErrInvalidAccessToken           = errors.New(ErrInvalidTokenHashText)
	ErrInvalidAccessTokenInBody     = errors.New(ErrInvalidAccessTokenInBodyText)
	ErrInvalidRefreshTokenInBody    = errors.New(ErrInvalidRefreshTokenInBodyText)
	ErrAccountNotFound              = errors.New(ErrAccountNotFoundText)
	ErrInvalidDeviceIDInBody        = errors.New(ErrInvalidDeviceIDInBodyText)
)
