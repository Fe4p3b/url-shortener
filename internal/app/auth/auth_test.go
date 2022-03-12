// Package auth provides business logic to authentication
// and authorization.
package auth

import (
	"crypto/cipher"
	"reflect"
	"testing"

	"github.com/Fe4p3b/url-shortener/internal/repositories"
)

func TestAuth_Encrypt(t *testing.T) {
	type fields struct {
		r      repositories.AuthRepository
		key    [32]byte
		aesgcm cipher.AEAD
	}
	type args struct {
		src string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Auth{
				r:      tt.fields.r,
				key:    tt.fields.key,
				aesgcm: tt.fields.aesgcm,
			}
			got, err := a.Encrypt(tt.args.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("Auth.Encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Auth.Encrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuth_Decrypt(t *testing.T) {
	type fields struct {
		r      repositories.AuthRepository
		key    [32]byte
		aesgcm cipher.AEAD
	}
	type args struct {
		src string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Auth{
				r:      tt.fields.r,
				key:    tt.fields.key,
				aesgcm: tt.fields.aesgcm,
			}
			got, err := a.Decrypt(tt.args.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("Auth.Decrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Auth.Decrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}
