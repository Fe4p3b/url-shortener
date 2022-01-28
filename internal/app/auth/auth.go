package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"
)

type AuthService interface {
	GenerateUUID() string
	Encrypt(string) (string, error)
	Decrypt(string) ([]byte, error)
}

type Auth struct {
	key []byte
}

var _ AuthService = &Auth{}

func NewAuth(key []byte) *Auth {
	return &Auth{key: key}
}

func (a *Auth) GenerateUUID() string {
	return uuid.NewString()
}

func (a *Auth) Encrypt(src string) (string, error) {
	key := sha256.Sum256([]byte(a.key))

	aesblock, err := aes.NewCipher(key[:])
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return "", err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return "", err
	}

	// создаём вектор инициализации
	nonce := key[len(key)-aesgcm.NonceSize():]
	dst := (aesgcm.Seal(nil, nonce, []byte(src), nil)) // зашифровываем
	return fmt.Sprintf("%x", dst), nil
}

func (a *Auth) Decrypt(src string) ([]byte, error) {
	key := sha256.Sum256([]byte(a.key))

	aesblock, err := aes.NewCipher(key[:])
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return nil, err
	}

	nonce := key[len(key)-aesgcm.NonceSize():]

	encrypted, err := hex.DecodeString(src)
	if err != nil {
		panic(err)
	}

	dst, err := aesgcm.Open(nil, nonce, encrypted, nil) // расшифровываем
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return nil, err
	}

	return dst, nil
}
