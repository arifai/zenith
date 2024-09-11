package crypto

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/crypto/argon2"
	"strings"
)

// Argon2IdHash is a struct that represents the argon2id hash
type Argon2IdHash struct {
	// Time is the time in seconds
	Time uint32

	// Memory is the memory in kilobytes
	Memory uint32

	// Threads is the threads used
	Threads uint8

	// KeyLen is the key length in bytes
	KeyLen uint32

	// SaltLen is the salt length in bytes
	SaltLen uint32
}

// GenerateHash will generate a hash and salt,
// the hash is generated using argon2id algorithm
func (a *Argon2IdHash) GenerateHash(password, salt []byte) (string, error) {
	var err error

	if len(salt) > 0 && uint32(len(salt)) != a.SaltLen {
		return "", fmt.Errorf("salt length is incorrect: expected %d bytes, got %d bytes", a.SaltLen, len(salt))
	}

	if len(salt) == 0 {
		salt, err = generateBytes(a.SaltLen)
		if err != nil {
			return "", err
		}
	}

	hash := argon2.IDKey(password, salt, a.Time, a.Memory, a.Threads, a.KeyLen)
	base64Salt := base64.StdEncoding.EncodeToString(salt)
	base64Hash := base64.StdEncoding.EncodeToString(hash)
	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, a.Memory, a.Time, a.Threads, base64Salt, base64Hash)

	return encodedHash, nil
}

// VerifyHash will verify the password with the encoded hash
func VerifyHash(password, encodedHash string) (bool, error) {
	a, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	otherHash := argon2.IDKey([]byte(password), salt, a.Time, a.Memory, a.Threads, a.KeyLen)

	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}

	return false, nil
}

// generateBytes will generate random bytes
func generateBytes(length uint32) ([]byte, error) {
	secret := make([]byte, length)
	_, err := rand.Read(secret)
	if err != nil {
		return nil, err
	}

	return secret, nil
}

// decodeHash will decode the hash
func decodeHash(encodedHash string) (a *Argon2IdHash, salt, hash []byte, err error) {
	value := strings.Split(encodedHash, "$")
	if len(value) != 6 {
		err = fmt.Errorf("invalid encoded hash")
		return nil, nil, nil, errors.New("invalid encoded hash")
	}

	var version int
	_, err = fmt.Sscanf(value[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, errors.New("incompatible Argon2 version")
	}

	a = &Argon2IdHash{}
	_, err = fmt.Sscanf(value[3], "m=%d,t=%d,p=%d", &a.Memory, &a.Time, &a.Threads)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.StdEncoding.DecodeString(value[4])
	if err != nil {
		return nil, nil, nil, err
	}

	a.SaltLen = uint32(len(salt))
	hash, err = base64.StdEncoding.DecodeString(value[5])
	if err != nil {
		return nil, nil, nil, err
	}

	a.KeyLen = uint32(len(hash))

	return a, salt, hash, nil
}
