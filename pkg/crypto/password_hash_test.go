package crypto

import (
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

type MockArgon2IdHash struct {
	mock.Mock
	Argon2IdHash
}

func (m *MockArgon2IdHash) GenerateHash(password, salt []byte) (string, error) {
	args := m.Called(password, salt)
	return args.String(0), args.Error(1)
}

func TestGenerateHash(t *testing.T) {
	password := []byte(faker.Password())
	salt := []byte(faker.Word())
	hashedPassword := faker.Password()

	mockArgon2IdHash := new(MockArgon2IdHash)
	mockArgon2IdHash.On("GenerateHash", password, salt).Return(hashedPassword, nil)

	result, err := mockArgon2IdHash.GenerateHash(password, salt)

	assert.NoError(t, err)
	assert.Equal(t, hashedPassword, result)

	mockArgon2IdHash.AssertExpectations(t)
}

func TestVerifyHash(t *testing.T) {
	password := faker.Password()

	argon2IdHash := &Argon2IdHash{
		Time:    1,
		Memory:  512,
		Threads: 2,
		KeyLen:  16,
		SaltLen: 16,
	}

	generatedHash, err := argon2IdHash.GenerateHash([]byte(password), nil)
	require.NoError(t, err, "Failed to generate hash")

	tests := []struct {
		name           string
		password       string
		hashedPassword string
		expectedResult bool
	}{
		{
			name:           "SuccessfulValidation",
			password:       password,
			hashedPassword: generatedHash,
			expectedResult: true,
		},
		{
			name:           "InvalidPassword",
			password:       faker.Password(),
			hashedPassword: generatedHash,
			expectedResult: false,
		},
		{
			name:           "InvalidHashFormat",
			password:       faker.Password(),
			hashedPassword: "$argon2i$v=19$invalid_hash$format",
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := VerifyHash(tt.password, tt.hashedPassword)
			assert.NoError(t, err, "Verification resulted in an error")
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}
