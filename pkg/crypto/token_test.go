package crypto

import (
	"aidanwoods.dev/go-paseto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

const InvalidToken = "InvalidToken"

type MockKey struct {
	mock.Mock
}

func (m *MockKey) ExportHex() paseto.V4AsymmetricSecretKey {
	args := m.Called()
	return args.Get(0).(paseto.V4AsymmetricSecretKey)
}

func (m *MockKey) ExportPublicKey() paseto.V4AsymmetricPublicKey {
	args := m.Called()
	return args.Get(0).(paseto.V4AsymmetricPublicKey)
}

func setupMockKey() *MockKey {
	mockKey := new(MockKey)
	mockedKey := paseto.NewV4AsymmetricSecretKey()
	mockedPublicKey := mockedKey.Public()
	mockKey.On("ExportHex").Return(mockedKey)
	mockKey.On("ExportPublicKey").Return(mockedPublicKey)

	return mockKey
}

func TestGenerateToken(t *testing.T) {
	mockKey := setupMockKey()
	jti := uuid.New()
	accountId := uuid.New()
	tokenPayload := TokenPayload{
		Jti:       jti,
		AccountId: accountId,
		IssuedAt:  time.Now(),
		NotBefore: time.Now(),
		ExpiresAt: time.Now().Add(time.Hour * 1),
		TokenType: "access",
	}
	token := tokenPayload.GenerateToken(mockKey.ExportHex())
	assert.NotEmpty(t, token)
}

func TestVerifyToken(t *testing.T) {
	mockKey := setupMockKey()
	validPayload := TokenPayload{
		Jti:       uuid.New(),
		AccountId: uuid.New(),
		IssuedAt:  time.Now(),
		NotBefore: time.Now(),
		ExpiresAt: time.Now().Add(time.Hour * 1),
		TokenType: "access",
	}
	validToken := validPayload.GenerateToken(mockKey.ExportHex())

	testCases := []struct {
		name        string
		token       string
		expectError bool
	}{
		{
			name:        "ValidTokenTest",
			token:       validToken,
			expectError: false,
		},
		{
			name:        "InvalidTokenTest",
			token:       InvalidToken,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer mockKey.AssertExpectations(t)
			_, err := VerifyToken(tc.token, mockKey.ExportPublicKey())
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestParseUUID(t *testing.T) {
	testCases := []struct {
		name        string
		input       func() (string, error)
		expectError bool
	}{
		{
			name:        "NormalCase",
			input:       func() (string, error) { return uuid.New().String(), nil },
			expectError: false,
		},
		{
			name:        "InvalidCase",
			input:       func() (string, error) { return "InvalidUUID", nil },
			expectError: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := parseUUID(tc.input, "Failed to get field", "Failed to parse field")
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
