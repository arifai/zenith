package crypto

import (
	"aidanwoods.dev/go-paseto"
	"errors"
	errmsg "github.com/arifai/go-modular-monolithic/internal/errors"
	"github.com/google/uuid"
	"log"
	"time"
)

// TokenPayload represents the payload data structure embedded in a token.
type TokenPayload struct {
	Jti       uuid.UUID
	AccountId uuid.UUID
	IssuedAt  time.Time
	NotBefore time.Time
	ExpiresAt time.Time
	TokenType string
}

// GenerateToken creates a signed token using the given secret key.
func (t *TokenPayload) GenerateToken(secretKey paseto.V4AsymmetricSecretKey) string {
	token := paseto.NewToken()

	token.SetJti(t.Jti.String())
	token.SetSubject(t.AccountId.String())
	token.SetIssuedAt(t.IssuedAt)
	token.SetNotBefore(t.NotBefore)
	token.SetExpiration(t.ExpiresAt)
	token.SetFooter([]byte(t.TokenType))

	signed := token.V4Sign(secretKey, nil)

	return signed
}

// VerifyToken verifies a given token using the provided public key, and returns the decoded TokenPayload if valid.
func VerifyToken(token string, publicKey paseto.V4AsymmetricPublicKey) (*TokenPayload, error) {
	parser := paseto.NewParser()
	parser.AddRule(paseto.NotBeforeNbf())
	parser.AddRule(paseto.ValidAt(time.Now()))

	pubKeyHex, err := paseto.NewV4AsymmetricPublicKeyFromHex(publicKey.ExportHex())
	if err != nil {
		log.Printf("%s: %v", errmsg.ErrFailedParsePublicHexText, err)
		return nil, err
	}

	parsed, err := parser.ParseV4Public(pubKeyHex, token, nil)
	if err != nil {
		log.Printf("%s: %v", errmsg.ErrFailedParseTokenText, err)
		return nil, errors.New(errmsg.ErrInvalidAccessTokenText)
	}

	jtiString, err := parsed.GetJti()
	if err != nil {
		log.Printf("%s: %v", errmsg.ErrFailedGetJTIText, err)
		return nil, err
	}

	jti, err := uuid.Parse(jtiString)
	if err != nil {
		log.Printf("%s: %v", errmsg.ErrFailedParseJTIText, err)
		return nil, err
	}

	accountIdString, err := parsed.GetSubject()
	if err != nil {
		log.Printf("%s: %v", errmsg.ErrFailedGetSubText, err)
		return nil, err
	}

	accountId, err := uuid.Parse(accountIdString)
	if err != nil {
		log.Printf("%s: %v", errmsg.ErrFailedParseACIText, err)
		return nil, err
	}

	issuedAt, err := parsed.GetIssuedAt()
	if err != nil {
		log.Printf("%s: %v", errmsg.ErrFailedGetIATText, err)
		return nil, err
	}

	notBefore, err := parsed.GetNotBefore()
	if err != nil {
		log.Printf("%s: %v", errmsg.ErrFailedGetNBFText, err)
		return nil, err
	}

	exp, err := parsed.GetExpiration()
	if err != nil {
		log.Printf("%s: %v", errmsg.ErrFailedGetEXPText, err)
		return nil, err
	} else if exp.Before(time.Now()) {
		log.Printf("%s: %v", errmsg.ErrTokenExpiredText, err)
		return nil, err
	}

	tokenPayload := &TokenPayload{
		Jti:       jti,
		AccountId: accountId,
		IssuedAt:  issuedAt,
		NotBefore: notBefore,
		ExpiresAt: exp,
		TokenType: string(parsed.Footer()),
	}

	return tokenPayload, nil
}
