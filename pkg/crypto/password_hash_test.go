package crypto

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockArgon2IdHash struct {
	mock.Mock
}

func (m *MockArgon2IdHash) GenerateHash(password, salt []byte) (string, error) {
	args := m.Called(password, salt)
	return args.String(0), args.Error(1)
}

var password = []byte("password")
var salt = []byte("salt")

func TestMockArgon2IdHash_GenerateHash(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		password     []byte
		salt         []byte
		expectedHash string
		expectedErr  error
		mockSetup    func(m *MockArgon2IdHash)
	}{
		{
			name:         "ValidHash",
			password:     password,
			salt:         salt,
			expectedHash: "HashedPassword",
			expectedErr:  nil,
			mockSetup: func(m *MockArgon2IdHash) {
				m.On("GenerateHash", password, salt).Return("HashedPassword", nil).Once()
			},
		},
		{
			name:         "ErrorGeneratingHash",
			password:     password,
			salt:         salt,
			expectedHash: "",
			expectedErr:  errors.New("error generating hash"),
			mockSetup: func(m *MockArgon2IdHash) {
				m.On("GenerateHash", password, salt).Return("", errors.New("error generating hash")).Once()
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mockArgon2IdHash := new(MockArgon2IdHash)
			tc.mockSetup(mockArgon2IdHash)

			hash, err := mockArgon2IdHash.GenerateHash(tc.password, tc.salt)

			assert.Equal(t, tc.expectedHash, hash)
			assert.Equal(t, tc.expectedErr, err)

			mockArgon2IdHash.AssertExpectations(t)
		})
	}
}
