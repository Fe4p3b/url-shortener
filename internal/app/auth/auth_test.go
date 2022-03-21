// Package auth provides business logic to authentication
// and authorization.
package auth

import (
	"log"
	"testing"

	"github.com/Fe4p3b/url-shortener/internal/storage/memory"
	"github.com/stretchr/testify/assert"
)

func TestAuth_Encrypt(t *testing.T) {
	m := memory.NewMemory(
		map[string]string{
			"asdf": "yandex.ru",
		},
	)
	auth, err := NewAuth([]byte("x35k9f"), m)
	if err != nil {
		log.Fatal(err)
	}

	tests := []struct {
		name    string
		auth    *Auth
		src     string
		want    string
		wantErr bool
	}{
		{
			name: "Test case #1",
			auth: auth,
			src:  "asdf",
			want: "ba75786072950c920d0202a7404cc47e317eb440",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := tt.auth.Encrypt(tt.src)
			log.Println(got)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAuth_Decrypt(t *testing.T) {
	m := memory.NewMemory(
		map[string]string{
			"asdf": "yandex.ru",
		},
	)
	auth, err := NewAuth([]byte("x35k9f"), m)
	if err != nil {
		log.Fatal(err)
	}
	tests := []struct {
		name    string
		auth    *Auth
		src     string
		want    string
		wantErr bool
	}{
		{
			name: "Test case #1",
			auth: auth,
			src:  "ba75786072950c920d0202a7404cc47e317eb440",
			want: "asdf",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.auth.Decrypt(tt.src)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, string(got))
		})
	}
}
