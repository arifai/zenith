package crypto_test

import (
	"github.com/arifai/go-modular-monolithic/pkg/crypto"
	"testing"
)

var (
	mockHash = crypto.Argon2IdHash{Time: 1, Memory: 64 * 1024, Threads: 4, KeyLen: 32, SaltLen: 16}
)

func TestGenerateHash(t *testing.T) {
	password := []byte("12345678")
	var salt []byte

	t.Run("SuccessNoSaltProvided", func(t *testing.T) {
		hash, err := mockHash.GenerateHash(password, salt)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if hash == "" {
			t.Fatalf("Expected non-empty hash, got empty string")
		}
	})

	t.Run("SuccessWithValidSalt", func(t *testing.T) {
		salt = []byte("1234567890123456")
		hash, err := mockHash.GenerateHash(password, salt)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if hash == "" {
			t.Fatalf("Expected non-empty hash, got empty string")
		}
	})

	t.Run("FailureWithShortSalt", func(t *testing.T) {
		salt = []byte("short")
		_, err := mockHash.GenerateHash(password, salt)
		if err == nil {
			t.Fatalf("Expected error due to short salt, but got none")
		}
	})
}

func TestVerify(t *testing.T) {
	password := []byte("12345678")
	wrongPassword := []byte("wrong_password")

	t.Run("Success", func(t *testing.T) {
		hash, err := mockHash.GenerateHash(password, nil)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		isValid, err := crypto.VerifyHash(string(password), hash)
		if err != nil {
			t.Fatalf("Expected no error during verification, got %v", err)
		}

		if !isValid {
			t.Fatalf("Expected password to match, but verification failed")
		}
	})

	t.Run("FailureWrongPassword", func(t *testing.T) {
		hash, err := mockHash.GenerateHash(password, nil)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		isValid, err := crypto.VerifyHash(string(wrongPassword), hash)
		if err != nil {
			t.Fatalf("Expected no error during verification, got %v", err)
		}

		if isValid {
			t.Fatalf("Expected verification to fail for wrong password, but it succeeded")
		}
	})

	t.Run("FailureInvalidHashFormat", func(t *testing.T) {
		invalidHash := "invalid_hash"

		_, err := crypto.VerifyHash(string(password), invalidHash)
		if err == nil {
			t.Fatalf("Expected error for invalid hash format, but got none")
		}
	})
}
