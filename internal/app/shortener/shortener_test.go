package shortener

import (
	"testing"

	"github.com/Fe4p3b/url-shortener/internal/repositories"
	"github.com/Fe4p3b/url-shortener/internal/storage/memory"
	"github.com/stretchr/testify/assert"
)

func Test_shortener_Find(t *testing.T) {
	type fields struct {
		r repositories.ShortenerRepository
	}

	storage := memory.NewMemory(
		map[string]string{
			"asdf": "yandex.ru",
		},
	)

	tests := []struct {
		name     string
		fields   fields
		shortURL string
		want     string
		wantErr  bool
	}{
		{
			name: "test case #1",
			fields: fields{
				r: storage,
			},
			shortURL: "asdf",
			want:     "yandex.ru",
			wantErr:  false,
		},
		{
			name: "test case #2",
			fields: fields{
				r: storage,
			},
			shortURL: "qwerty",
			want:     "",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &shortener{
				r: tt.fields.r,
			}
			got, err := s.Find(tt.shortURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("shortener.Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_shortener_Store(t *testing.T) {
	type fields struct {
		r repositories.ShortenerRepository
	}

	storage := memory.NewMemory(
		map[string]string{
			"asdf": "yandex.ru",
		},
	)
	tests := []struct {
		name    string
		fields  fields
		urls    []string
		wantErr bool
	}{
		{
			name: "test for duplicate short urls",
			fields: fields{
				r: storage,
			},
			urls:    []string{"google.com", "yandex.ru", "yahoo.com", "google.com", "yandex.ru", "yahoo.com"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &shortener{
				r: tt.fields.r,
			}
			shortURLs := make(map[string]int)
			for _, v := range tt.urls {
				got, err := s.Store(v)
				if (err != nil) != tt.wantErr {
					t.Errorf("shortener.Store() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				shortURLs[got]++
			}
			t.Log(shortURLs)
			for k, v := range shortURLs {
				if v > 1 {
					t.Errorf("shortener.Store() duplicate short urls generated %v, in %v", k, s.r)
					return
				}
			}

		})
	}
}
