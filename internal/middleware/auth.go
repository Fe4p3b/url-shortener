// Package middleware provides functionality for wrapper functions
// for handlers, that perform required operations either before
// calling handler or after.
package middleware

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/Fe4p3b/url-shortener/internal/app/auth"
)

type ContextKey string

var Key ContextKey = "user"

// AuthMiddleware is a middleware for authentication.
// If user provides token in a cookie, named token, it
// tries to authenticate him. If user can't be authenticated
// new token is created for user, and passed in a cookie,
// named token.
type AuthMiddleware struct {
	// auth is a service that performs operations on
	// encryption, decryption and authentication.
	auth auth.AuthService
}

func NewAuthMiddleware(auth auth.AuthService) *AuthMiddleware {
	return &AuthMiddleware{auth: auth}
}

// Middleware is a function that is used to wrap handler.
func (a *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("token")
		if err == nil {
			user, err := a.auth.Decrypt(token.Value)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			err = a.auth.VerifyUser(string(user))
			if err == nil {
				ctx := context.WithValue(r.Context(), Key, string(user))
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			if !errors.Is(err, sql.ErrNoRows) {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}

		uuid, err := a.auth.CreateUser()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		user, err := a.auth.Encrypt(uuid)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{Name: "token", Value: user})
		ctx := context.WithValue(r.Context(), Key, uuid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
