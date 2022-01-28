package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Fe4p3b/url-shortener/internal/app/auth"
)

type ContextKey string

var Key ContextKey = "user"

type AuthMiddleware struct {
	auth auth.AuthService
}

func NewAuthMiddleware(auth auth.AuthService) *AuthMiddleware {
	return &AuthMiddleware{auth: auth}
}

func (a *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("token")
		if err != nil {
			uuid := a.auth.GenerateUUID()

			user, err := a.auth.Encrypt(uuid)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			fmt.Printf("encrypted %s\n", user)

			http.SetCookie(w, &http.Cookie{Name: "token", Value: user})
			ctx := context.WithValue(r.Context(), Key, uuid)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		user, err := a.auth.Decrypt(token.Value)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		fmt.Printf("decrypted: %s\n", user)

		ctx := context.WithValue(r.Context(), Key, string(user))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
