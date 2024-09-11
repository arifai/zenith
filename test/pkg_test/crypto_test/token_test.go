package crypto

import (
	"aidanwoods.dev/go-paseto"
	crp "github.com/arifai/go-modular-monolithic/pkg/crypto"
	"github.com/google/uuid"
	"testing"
	"time"
)

var (
	u         = uuid.New()
	tn        = time.Now()
	mockToken = crp.TokenPayload{
		Jti:       u,
		AccountId: u,
		IssuedAt:  tn,
		NotBefore: tn,
		ExpiresAt: tn.Add(time.Hour * 24),
	}
)

func timeAlmostEqual(t1, t2 time.Time, tolerance time.Duration) bool {
	return t1.Sub(t2) < tolerance && t2.Sub(t1) < tolerance
}

func TestGenerateToken(t *testing.T) {
	secretKey := paseto.NewV4AsymmetricSecretKey()

	t.Run("Success", func(t *testing.T) {
		token := mockToken.GenerateToken(secretKey)
		if token == "" {
			t.Fatalf("Expected a valid token, got an empty string")
		}

		publicKey := secretKey.Public()
		v4 := paseto.NewParserWithoutExpiryCheck()
		parsedToken, err := v4.ParseV4Public(publicKey, token, nil)
		if err != nil {
			t.Fatalf("Failed to parse token: %v", err)
		}

		jti, err := parsedToken.GetJti()
		if err != nil {
			t.Fatalf("Failed to get 'jti' claim: %v", err)
		}

		if jti != u.String() {
			t.Fatalf("Expected 'jti' %v, got %v", u.String(), jti)
		}

		subject, err := parsedToken.GetSubject()
		if err != nil {
			t.Fatalf("Failed to get 'subject' claim: %v", err)
		}

		if subject != u.String() {
			t.Fatalf("Expected 'subject' %v, got %v", u.String(), subject)
		}
	})

	t.Run("ExpiredToken", func(t *testing.T) {
		token := mockToken.GenerateToken(secretKey)
		publicKey := secretKey.Public()
		v4 := paseto.NewParser()

		parsedToken, err := v4.ParseV4Public(publicKey, token, nil)
		if err != nil {
			t.Fatalf("Failed to parse token: %v", err)
		}

		expiry, err := parsedToken.GetExpiration()
		if err != nil {
			t.Fatalf("Failed to get expiration time: %v", err)
		}

		if !expiry.After(tn) {
			t.Fatal("Expected token not to be expired immediately")
		}
	})
}

func TestVerifyToken(t *testing.T) {
	secretKey := paseto.NewV4AsymmetricSecretKey()
	publicKey := secretKey.Public()
	token := mockToken.GenerateToken(secretKey)
	tolerance := time.Second

	t.Run("ValidToken", func(t *testing.T) {
		parsedPayload, err := crp.VerifyToken(token, publicKey)
		if err != nil {
			t.Fatalf("Failed to verify valid token: %v", err)
		}

		if parsedPayload.Jti != mockToken.Jti {
			t.Errorf("Expected Jti %v, got %v", mockToken.Jti, parsedPayload.Jti)
		}

		if parsedPayload.AccountId != mockToken.AccountId {
			t.Errorf("Expected AccountId %v, got %v", mockToken.AccountId, parsedPayload.AccountId)
		}

		if !timeAlmostEqual(parsedPayload.IssuedAt, mockToken.IssuedAt, tolerance) {
			t.Errorf("Expected IssuedAt %v, got %v", mockToken.IssuedAt, parsedPayload.IssuedAt)
		}

		if !timeAlmostEqual(parsedPayload.NotBefore, mockToken.NotBefore, tolerance) {
			t.Errorf("Expected NotBefore %v, got %v", mockToken.NotBefore, parsedPayload.NotBefore)
		}

		if !timeAlmostEqual(parsedPayload.ExpiresAt, mockToken.ExpiresAt, tolerance) {
			t.Errorf("Expected ExpiresAt %v, got %v", mockToken.ExpiresAt, parsedPayload.ExpiresAt)
		}
	})

	t.Run("ExpiredToken", func(t *testing.T) {
		expiredPayload := &crp.TokenPayload{
			Jti:       u,
			AccountId: u,
			IssuedAt:  tn.Add(-48 * time.Hour),
			NotBefore: tn.Add(-48 * time.Hour),
			ExpiresAt: tn.Add(-24 * time.Hour),
		}
		expiredToken := expiredPayload.GenerateToken(secretKey)

		_, err := crp.VerifyToken(expiredToken, publicKey)
		if err == nil {
			t.Fatal("Expected an error when verifying expired token, but got none")
		}
	})

	t.Run("InvalidSignature", func(t *testing.T) {
		invalidSecretKey := paseto.NewV4AsymmetricSecretKey()
		invalidToken := mockToken.GenerateToken(invalidSecretKey)

		_, err := crp.VerifyToken(invalidToken, publicKey)
		if err == nil {
			t.Fatal("Expected an error when verifying a token with an invalid signature, but got none")
		}
	})

	t.Run("MalformedToken", func(t *testing.T) {
		malformedToken := "this.is.not.a.valid.token"

		_, err := crp.VerifyToken(malformedToken, publicKey)
		if err == nil {
			t.Fatal("Expected an error when verifying malformed token, but got none")
		}
	})
}
