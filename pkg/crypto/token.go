package crypto

import (
	"aidanwoods.dev/go-paseto"
	"fmt"
	"github.com/google/uuid"
	"time"
)

// TokenPayload struct to hold token payload
type TokenPayload struct {
	// Jti is the unique identifier of the token
	Jti uuid.UUID

	// AccountId is the unique identifier of the account. This AccountId will insert into the 'sub' field
	AccountId uuid.UUID

	// IssuedAt is the time when the token is issued
	IssuedAt time.Time

	// NotBefore is the time when the token is valid
	NotBefore time.Time

	// ExpiresAt is the time when the token is expired
	ExpiresAt time.Time
}

// GenerateToken function to generate token
func (t *TokenPayload) GenerateToken(secretKey paseto.V4AsymmetricSecretKey) string {
	token := paseto.NewToken()

	token.SetJti(t.Jti.String())
	token.SetSubject(t.AccountId.String())
	token.SetIssuedAt(t.IssuedAt)
	token.SetNotBefore(t.NotBefore)
	token.SetExpiration(t.ExpiresAt)

	signed := token.V4Sign(secretKey, nil)

	return signed
}

// VerifyToken function to verify token
func VerifyToken(token string, publicKey paseto.V4AsymmetricPublicKey) (*TokenPayload, error) {
	parser := paseto.NewParser()
	parser.AddRule(paseto.NotBeforeNbf())
	parser.AddRule(paseto.ValidAt(time.Now()))

	parsed, err := parser.ParseV4Public(publicKey, token, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	jtiString, err := parsed.GetJti()
	if err != nil {
		return nil, fmt.Errorf("failed to get 'jti': %w", err)
	}

	jti, err := uuid.Parse(jtiString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse jti string: %w", err)
	}

	accountIdString, err := parsed.GetSubject()
	if err != nil {
		return nil, fmt.Errorf("failed to get 'sub': %w", err)
	}

	accountId, err := uuid.Parse(accountIdString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse account id string: %w", err)
	}

	issuedAt, err := parsed.GetIssuedAt()
	if err != nil {
		return nil, fmt.Errorf("failed to get 'iat': %w", err)
	}

	notBefore, err := parsed.GetNotBefore()
	if err != nil {
		return nil, fmt.Errorf("failed to get 'nbf': %w", err)
	}

	exp, err := parsed.GetExpiration()
	if err != nil {
		return nil, fmt.Errorf("failed to get 'exp': %w", err)
	}

	tokenPayload := &TokenPayload{Jti: jti, AccountId: accountId, IssuedAt: issuedAt, NotBefore: notBefore, ExpiresAt: exp}

	return tokenPayload, nil
}
