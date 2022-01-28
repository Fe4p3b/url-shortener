package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/Fe4p3b/url-shortener/internal/repositories"
)

type AuthService interface {
	CreateUser() (string, error)
	Encrypt(string) (string, error)
	Decrypt(string) ([]byte, error)
	VerifyUser(string) error
}

type Auth struct {
	r   repositories.AuthRepository
	key []byte
}

var _ AuthService = &Auth{}

func NewAuth(key []byte, r repositories.AuthRepository) *Auth {
	return &Auth{key: key, r: r}
}

func (a *Auth) CreateUser() (string, error) {
	return a.r.CreateUser()
}

func (a *Auth) Encrypt(src string) (string, error) {
	key := sha256.Sum256([]byte(a.key))

	aesblock, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return "", err
	}

	nonce := key[len(key)-aesgcm.NonceSize():]
	dst := (aesgcm.Seal(nil, nonce, []byte(src), nil))
	return fmt.Sprintf("%x", dst), nil
}

func (a *Auth) Decrypt(src string) ([]byte, error) {
	key := sha256.Sum256([]byte(a.key))

	aesblock, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, err
	}

	nonce := key[len(key)-aesgcm.NonceSize():]

	encrypted, err := hex.DecodeString(src)
	if err != nil {
		panic(err)
	}

	dst, err := aesgcm.Open(nil, nonce, encrypted, nil)
	if err != nil {
		return nil, err
	}

	return dst, nil
}

func (a *Auth) VerifyUser(user string) error {
	return a.r.VerifyUser(user)
}
