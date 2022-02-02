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
	key [32]byte
}

var _ AuthService = &Auth{}

func NewAuth(key []byte, r repositories.AuthRepository) *Auth {
	return &Auth{key: sha256.Sum256([]byte(key)), r: r}
}

func (a *Auth) CreateUser() (string, error) {
	return a.r.CreateUser()
}

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

func (a *Auth) Encrypt(src string) (string, error) {
	aesgcm, err := a.GetAesgcm()
	if err != nil {
		return "", err
	}

	nonce := a.key[len(a.key)-aesgcm.NonceSize():]
	dst := (aesgcm.Seal(nil, nonce, []byte(src), nil))
	return fmt.Sprintf("%x", dst), nil
}

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

func (a *Auth) VerifyUser(user string) error {
	return a.r.VerifyUser(user)
}
