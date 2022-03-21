// Package auth provides business logic to authentication
// and authorization.
package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/Fe4p3b/url-shortener/internal/repositories"
)

// AuthService service for authentication and authorization.
type AuthService interface {
	// CreateUser creates user or returns error.
	CreateUser() (string, error)

	// Encrypt encrypts string using GCM and AES256 algorithms and
	// returns encrypted string, or error.
	Encrypt(string) (string, error)

	// Decrypt decrypts string using GCM and AES256 algorithms and
	// returns slice of bytes, or error
	Decrypt(string) ([]byte, error)

	// VerifyUser authenticates user, by encrypted string, or returns
	// error
	VerifyUser(string) error
}

// Auth provides functionality for authentication and authorization.
type Auth struct {
	// r is a storage.
	r repositories.AuthRepository

	// key is a 32-byte key, that is used for encryption and
	// decryption.
	key [32]byte

	aesgcm cipher.AEAD
}

var _ AuthService = &Auth{}

func NewAuth(key []byte, r repositories.AuthRepository) (*Auth, error) {
	authKey := sha256.Sum256([]byte(key))
	aesblock, err := aes.NewCipher(authKey[:])
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, err
	}

	return &Auth{key: authKey, r: r, aesgcm: aesgcm}, nil
}

// CreateUser implements ShortenerService CreateUser method.
func (a *Auth) CreateUser() (string, error) {
	return a.r.CreateUser()
}

// Encrypt implements ShortenerService Encrypt method.
func (a *Auth) Encrypt(src string) (string, error) {
	nonce := a.key[len(a.key)-a.aesgcm.NonceSize():]
	dst := (a.aesgcm.Seal(nil, nonce, []byte(src), nil))

	return fmt.Sprintf("%x", dst), nil
}

// Decrypt implements ShortenerService Decrypt method.
func (a *Auth) Decrypt(src string) ([]byte, error) {
	nonce := a.key[len(a.key)-a.aesgcm.NonceSize():]
	encrypted, err := hex.DecodeString(src)
	if err != nil {
		return nil, err
	}

	dst, err := a.aesgcm.Open(nil, nonce, encrypted, nil)
	if err != nil {
		return nil, err
	}

	return dst, nil
}

// VerifyUser implements ShortenerService VerifyUser method.
func (a *Auth) VerifyUser(user string) error {
	return a.r.VerifyUser(user)
}
