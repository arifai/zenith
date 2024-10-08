package crypto

import (
	"aidanwoods.dev/go-paseto"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"github.com/arifai/zenith/pkg/errormessage"
	"github.com/google/uuid"
	"io"
	"log"
	"time"
)

// TokenPayload represents the payload data structure embedded in a token.
type TokenPayload struct {
	Jti       uuid.UUID
	AccountId uuid.UUID
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
	token.SetSubject(t.AccountId.String())
	token.SetIssuedAt(t.IssuedAt)
	token.SetNotBefore(t.NotBefore)
	token.SetExpiration(t.ExpiresAt)
	token.SetFooter([]byte(t.TokenType))

	return token.V4Sign(secretKey, nil)
}

// GenerateHMAC takes a message string and generates the HMAC using a random secret key.
func GenerateHMAC(message string) (string, error) {
	randomSecret, err := generateRandomSecret(64)
	if err != nil {
		return "", err
	}

	key := []byte(randomSecret)
	h := hmac.New(sha512.New, key)
	h.Write([]byte(message))

	return hex.EncodeToString(h.Sum(nil)), nil
}

// generateRandomSecret generates a random secret key of the given length
func generateRandomSecret(length int) (string, error) {
	key := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return "", err
	}

	return hex.EncodeToString(key), nil
}

// VerifyToken verifies a given token using the provided public key, and returns the decoded TokenPayload if valid.
func VerifyToken(token string, publicKey paseto.V4AsymmetricPublicKey) (*TokenPayload, error) {
	parser := paseto.NewParser()
	parser.AddRule(paseto.NotBeforeNbf())
	parser.AddRule(paseto.ValidAt(time.Now()))

	publicKeyHex, err := paseto.NewV4AsymmetricPublicKeyFromHex(publicKey.ExportHex())
	if err != nil {
		log.Printf("%s: %v", errormessage.ErrFailedParsePublicHexText, err)
		return nil, err
	}

	parsedToken, err := parser.ParseV4Public(publicKeyHex, token, nil)
	if err != nil {
		log.Printf("%s: %v", errormessage.ErrFailedParseTokenText, err)
		return nil, errors.New(errormessage.ErrInvalidTokenHashText)
	}

	jti, err := parseUUID(parsedToken.GetJti, errormessage.ErrFailedGetJTIText, errormessage.ErrFailedParseJTIText)
	if err != nil {
		return nil, err
	}

	aud, err := parseUUID(parsedToken.GetAudience, errormessage.ErrFailedGetJTIText, errormessage.ErrFailedParseJTIText)
	if err != nil {
		return nil, err
	}

	accountId, err := parseUUID(parsedToken.GetSubject, errormessage.ErrFailedGetSubText, errormessage.ErrFailedParseACIText)
	if err != nil {
		return nil, err
	}

	issuedAt, err := parsedToken.GetIssuedAt()
	if err != nil {
		log.Printf("%s: %v", errormessage.ErrFailedGetIATText, err)
		return nil, err
	}

	notBefore, err := parsedToken.GetNotBefore()
	if err != nil {
		log.Printf("%s: %v", errormessage.ErrFailedGetNBFText, err)
		return nil, err
	}

	expiration, err := parsedToken.GetExpiration()
	if err != nil {
		log.Printf("%s: %v", errormessage.ErrFailedGetEXPText, err)
		return nil, err
	} else if expiration.Before(time.Now()) {
		log.Printf("%s: %v", errormessage.ErrTokenExpiredText, err)
		return nil, err
	}

	tokenPayload := &TokenPayload{
		Jti:       jti,
		DeviceID:  aud,
		AccountId: accountId,
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
		log.Printf("%s: %v", getFieldErrMsg, err)
		return uuid.Nil, err
	}

	fieldUUID, err := uuid.Parse(fieldStr)
	if err != nil {
		log.Printf("%s: %v", parseErrMsg, err)
		return uuid.Nil, err
	}

	return fieldUUID, nil
}
