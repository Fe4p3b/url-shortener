package middleware

import (
	"context"
	"crypto/aes"
	"log"
	"net/http"

	"github.com/teris-io/shortid"
)

type ContextKey string

var Key ContextKey = "user"

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("token")
		if err != nil {
			// uuid := uuid.NewString()
			uuid, err := shortid.Generate()
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			// b, err := Generate(uuid)
			// if err != nil {
			// 	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			// 	return
			// }

			log.Printf("uuid - %v, token - %s", uuid, "")
			http.SetCookie(w, &http.Cookie{Name: "token", Value: string(uuid)})
			next.ServeHTTP(w, r)
			return
		}

		// log.Printf("token - %v", token)

		// key := []byte("qwertyuioasdfghjqwertyuioasdfghj")

		// aesblock, err := aes.NewCipher(key)
		// if err != nil {
		// 	fmt.Printf("error: %v\n", err)
		// 	return
		// }

		// fmt.Printf("encrypted: %s\n", token.Value)
		// src := make([]byte, len(key)) // расшифровываем
		// aesblock.Decrypt(src, []byte(token.Value))
		// fmt.Printf("decrypted: %s\n", src)

		ctx := context.WithValue(r.Context(), Key, token.Value)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Generate(src string) ([]byte, error) {
	key := []byte("qwertyuioasdfghj")
	// qwert yuioa sdfgh jqwer tyuio asdfghj
	aesblock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	log.Printf("src - %s", []byte(src))
	dst := make([]byte, len(key)) // зашифровываем
	aesblock.Encrypt(dst, []byte(src))
	return dst, nil
}
