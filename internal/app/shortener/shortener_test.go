package shortener

import (
	"testing"
	"time"

	"github.com/Fe4p3b/url-shortener/internal/models"
	"github.com/Fe4p3b/url-shortener/internal/repositories"
	"github.com/Fe4p3b/url-shortener/internal/storage"
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

			if got != nil {
				assert.Equal(t, tt.want, got.URL)
			}
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
		urls    []models.URL
		wantErr bool
	}{
		{
			name: "test for duplicate short urls",
			fields: fields{
				r: storage,
			},
			urls: []models.URL{
				{URL: "google.com"},
				{URL: "yandex.ru"},
				{URL: "yahoo.com"},
				{URL: "google.com"},
				{URL: "yandex.ru"},
				{URL: "yahoo.com"},
			},
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
				got, err := s.Store(&v)
				if (err != nil) != tt.wantErr {
					t.Errorf("shortener.Store() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				shortURLs[got]++
			}

			for k, v := range shortURLs {
				if v > 1 {
					t.Errorf("shortener.Store() duplicate short urls generated %v, in %v", k, s.r)
					return
				}
			}

		})
	}
}

func Test_shortener_Ping(t *testing.T) {
	s := memory.NewMemory(
		map[string]string{
			"asdf": "yandex.ru",
		},
	)
	tests := []struct {
		name    string
		r       repositories.ShortenerRepository
		wantErr bool
	}{
		{
			name:    "Test case #1",
			r:       s,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &shortener{
				r: tt.r,
			}
			err := s.Ping()
			assert.NoError(t, err)
		})
	}
}

func Test_shortener_GetUserURLs(t *testing.T) {
	s := memory.NewMemory(
		map[string]string{
			"asdf": "yandex.ru",
		},
	)
	tests := []struct {
		name    string
		r       repositories.ShortenerRepository
		wantErr bool
	}{
		{
			name:    "Test case #1",
			r:       s,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &shortener{
				r: tt.r,
			}
			_, err := s.GetUserURLs("user")
			assert.Error(t, storage.ErrorMethodIsNotImplemented, err)
		})
	}
}

func Test_shortener_StoreBatch(t *testing.T) {
	s := memory.NewMemory(
		map[string]string{
			"asdf": "yandex.ru",
		},
	)
	tests := []struct {
		name    string
		r       repositories.ShortenerRepository
		args    []repositories.URL
		wantErr bool
	}{
		{
			name:    "Test case #1",
			r:       s,
			args:    make([]repositories.URL, 0),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &shortener{
				r: tt.r,
			}
			_, err := s.StoreBatch("user", tt.args)
			assert.Error(t, storage.ErrorMethodIsNotImplemented, err)
		})
	}
}

func Test_shortener_DeleteURLs(t *testing.T) {
	s := memory.NewMemory(
		map[string]string{
			"asdf": "yandex.ru",
		},
	)
	tests := []struct {
		name    string
		r       repositories.ShortenerRepository
		args    []string
		wantErr bool
	}{
		{
			name:    "Test case #1",
			r:       s,
			args:    make([]string, 0),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &shortener{
				r: tt.r,
			}
			s.DeleteURLs("user", tt.args)
			time.Sleep(1 * time.Second)
		})
	}
}
