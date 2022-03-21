// Package middleware provides functionality for wrapper functions
// for handlers, that perform required operations either before
// calling handler or after.
package middleware

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Fe4p3b/url-shortener/internal/app/auth"
	"github.com/Fe4p3b/url-shortener/internal/storage/memory"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware_Middleware(t *testing.T) {
	m := memory.NewMemory(
		map[string]string{
			"asdf": "yandex.ru",
		},
	)
	a, err := auth.NewAuth([]byte("x35k9f"), m)
	if err != nil {
		log.Fatal(err)
	}

	type fields struct {
		auth auth.AuthService
	}
	type args struct {
		next  http.Handler
		token string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test case #1",
			fields: fields{
				auth: a,
			},
			args: args{
				next:  http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
				token: "ba75786072950c920d0202a7404cc47e317eb440",
			},
			wantErr: true,
		},
		{
			name: "Test case #2",
			fields: fields{
				auth: a,
			},
			args: args{
				next: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AuthMiddleware{
				auth: tt.fields.auth,
			}

			request := httptest.NewRequest(http.MethodGet, "/", nil)
			ctx := context.WithValue(context.Background(), Key, tt.args.token)
			w := httptest.NewRecorder()
			if tt.args.token != "" {
				request.AddCookie(&http.Cookie{Name: "token", Value: tt.args.token})
			}

			a.Middleware(tt.args.next).ServeHTTP(w, request.WithContext(ctx))
			_, ok := request.Context().Value(Key).(string)
			if tt.wantErr {
				assert.False(t, ok)
			}

		})
	}
}
