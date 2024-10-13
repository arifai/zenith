package crypto

import (
	"aidanwoods.dev/go-paseto"
	"errors"
	"github.com/arifai/zenith/pkg/errormessage"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

// TokenPayload represents the payload data structure embedded in a token.
type TokenPayload struct {
	Jti       uuid.UUID
	AccountID uuid.UUID
	DeviceID  uuid.UUID
	IssuedAt  time.Time
	NotBefore time.Time
	ExpiresAt time.Time
	TokenType string
}

const (
	AccessToken  = "access_token"
	RefreshToken = "refresh_token"
)

// GenerateToken creates a signed token using the given secret key.
func (t *TokenPayload) GenerateToken(secretKey paseto.V4AsymmetricSecretKey) string {
	token := paseto.NewToken()
	token.SetAudience(t.DeviceID.String())
	token.SetJti(t.Jti.String())
	token.SetSubject(t.AccountID.String())
	token.SetIssuedAt(t.IssuedAt)
	token.SetNotBefore(t.NotBefore)
	token.SetExpiration(t.ExpiresAt)
	token.SetFooter([]byte(t.TokenType))

	return token.V4Sign(secretKey, nil)
}

// VerifyToken verifies a given token using the provided public key, and returns the decoded TokenPayload if valid.
func VerifyToken(token string, publicKey paseto.V4AsymmetricPublicKey) (*TokenPayload, error) {
	parser := paseto.NewParser()
	parser.AddRule(paseto.NotBeforeNbf())
	parser.AddRule(paseto.ValidAt(time.Now()))

	publicKeyHex, err := paseto.NewV4AsymmetricPublicKeyFromHex(publicKey.ExportHex())
	if err != nil {
		log.Error(errormessage.ErrFailedParsePublicHexText, zap.Error(err))
		return nil, err
	}

	parsedToken, err := parser.ParseV4Public(publicKeyHex, token, nil)
	if err != nil {
		log.Error(errormessage.ErrFailedParseTokenText, zap.Error(err))
		return nil, errors.New(errormessage.ErrInvalidTokenHashText)
	}

	jti, err := parseUUID(parsedToken.GetJti, errormessage.ErrFailedGetJTIText, errormessage.ErrFailedParseJTIText)
	if err != nil {
		return nil, err
	}

	aud, err := parseUUID(parsedToken.GetAudience, errormessage.ErrFailedGetAudText, errormessage.ErrFailedParseAudText)
	if err != nil {
		return nil, err
	}

	accountId, err := parseUUID(parsedToken.GetSubject, errormessage.ErrFailedGetSubText, errormessage.ErrFailedParseACIText)
	if err != nil {
		return nil, err
	}

	issuedAt, err := parsedToken.GetIssuedAt()
	if err != nil {
		log.Error(errormessage.ErrFailedGetIATText, zap.Error(err))
		return nil, err
	}

	notBefore, err := parsedToken.GetNotBefore()
	if err != nil {
		log.Error(errormessage.ErrFailedGetNBFText, zap.Error(err))
		return nil, err
	}

	expiration, err := parsedToken.GetExpiration()
	if err != nil {
		log.Error(errormessage.ErrFailedGetEXPText, zap.Error(err))
		return nil, err
	} else if expiration.Before(time.Now()) {
		log.Error(errormessage.ErrTokenExpiredText, zap.Error(err))
		return nil, err
	}

	tokenPayload := &TokenPayload{
		Jti:       jti,
		DeviceID:  aud,
		AccountID: accountId,
		IssuedAt:  issuedAt,
		NotBefore: notBefore,
		ExpiresAt: expiration,
		TokenType: string(parsedToken.Footer()),
	}

	return tokenPayload, nil
}

// parseUUID attempts to get a field and parse it as a UUID.
// It uses getFieldFunc to retrieve the field value as a string.
// Logs and returns an error if retrieval or parsing fails, using getFieldErrMsg and parseErrMsg respectively.
func parseUUID(getFieldFunc func() (string, error), getFieldErrMsg, parseErrMsg string) (uuid.UUID, error) {
	fieldStr, err := getFieldFunc()
	if err != nil {
		log.Error(getFieldErrMsg, zap.Error(err))
		return uuid.Nil, err
	}

	fieldUUID, err := uuid.Parse(fieldStr)
	if err != nil {
		log.Error(parseErrMsg, zap.Error(err))
		return uuid.Nil, err
	}

	return fieldUUID, nil
}
