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
}

var _ AuthService = &Auth{}

func NewAuth(key []byte, r repositories.AuthRepository) *Auth {
	return &Auth{key: sha256.Sum256([]byte(key)), r: r}
}

// CreateUser implements ShortenerService CreateUser method.
func (a *Auth) CreateUser() (string, error) {
	return a.r.CreateUser()
}

// GetAesgcm uses AES and GCM algorithms to create cipher.AEAD,
// that is used for encryption and decryption.
func (a *Auth) GetAesgcm() (cipher.AEAD, error) {
	aesblock, err := aes.NewCipher(a.key[:])
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, err
	}

	return aesgcm, nil
}

// Encrypt implements ShortenerService Encrypt method.
func (a *Auth) Encrypt(src string) (string, error) {
	aesgcm, err := a.GetAesgcm()
	if err != nil {
		return "", err
	}

	nonce := a.key[len(a.key)-aesgcm.NonceSize():]
	dst := (aesgcm.Seal(nil, nonce, []byte(src), nil))
	return fmt.Sprintf("%x", dst), nil
}

// Decrypt implements ShortenerService Decrypt method.
func (a *Auth) Decrypt(src string) ([]byte, error) {
	aesgcm, err := a.GetAesgcm()
	if err != nil {
		return nil, err
	}

	nonce := a.key[len(a.key)-aesgcm.NonceSize():]
	encrypted, err := hex.DecodeString(src)
	if err != nil {
		return nil, err
	}

	dst, err := aesgcm.Open(nil, nonce, encrypted, nil)
	if err != nil {
		return nil, err
	}

	return dst, nil
}

// VerifyUser implements ShortenerService VerifyUser method.
func (a *Auth) VerifyUser(user string) error {
	return a.r.VerifyUser(user)
}
