package memory

import (
	"testing"

	"github.com/Fe4p3b/url-shortener/internal/models"
	"github.com/Fe4p3b/url-shortener/internal/repositories"
	"github.com/Fe4p3b/url-shortener/internal/storage"
	"github.com/stretchr/testify/assert"
)

func Test_memory_Find(t *testing.T) {
	s := &Memory{
		S: map[string]string{
			"asdf": "yandex.ru",
		},
	}
	tests := []struct {
		name    string
		value   string
		want    string
		wantErr error
	}{
		{
			name:    "test case #1",
			value:   "asdf",
			want:    "yandex.ru",
			wantErr: nil,
		},
		{
			name:    "test case #2",
			value:   "qwer",
			want:    "",
			wantErr: storage.ErrorNoLinkFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Find(tt.value)

			if err != nil && tt.wantErr != err {
				t.Errorf("Find() error = %v, wantErr = %v", err, tt.wantErr)
			}

			if got != nil {
				assert.Equal(t, tt.want, got.URL)
			}
		})
	}
}

func Test_memory_Save(t *testing.T) {
	s := &Memory{
		S: map[string]string{
			"asdf": "yandex.ru",
		},
	}
	type args struct {
		url   models.URL
		short string
	}

	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr error
	}{
		{
			name: "test case #1",
			args: args{
				url: models.URL{URL: "google.com", ShortURL: "qwerty"},
			},
			wantErr: nil,
		},
		{
			name: "test case #2",
			args: args{
				url: models.URL{URL: "yahoo.com", ShortURL: "asdf"},
			},
			wantErr: storage.ErrorDuplicateShortlink,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.Save(&tt.args.url)

			if tt.wantErr != nil {
				assert.Error(t, tt.wantErr, err)
			} else {
				assert.NoError(t, tt.wantErr, err)
			}
			_, ok := s.S[tt.args.url.ShortURL]
			assert.Equal(t, true, ok)

		})
	}
}

func TestMemory_Ping(t *testing.T) {
	s := &Memory{
		S: map[string]string{
			"asdf": "yandex.ru",
		},
	}
	err := s.Ping()
	assert.NoError(t, err)
}

func TestMemory_GetUserURLs(t *testing.T) {
	s := &Memory{
		S: map[string]string{
			"asdf": "yandex.ru",
		},
	}
	_, err := s.GetUserURLs("", "")
	assert.Error(t, storage.ErrorMethodIsNotImplemented, err)
}

func TestMemory_AddURLBuffer(t *testing.T) {
	s := &Memory{
		S: map[string]string{
			"asdf": "yandex.ru",
		},
	}
	err := s.AddURLBuffer(repositories.URL{})
	assert.Error(t, storage.ErrorMethodIsNotImplemented, err)
}

func TestMemory_Flush(t *testing.T) {
	s := &Memory{
		S: map[string]string{
			"asdf": "yandex.ru",
		},
	}
	err := s.Flush()
	assert.Error(t, storage.ErrorMethodIsNotImplemented, err)
}

func TestMemory_FlushToDelete(t *testing.T) {
	s := &Memory{
		S: map[string]string{
			"asdf": "yandex.ru",
		},
	}
	err := s.FlushToDelete()
	assert.Error(t, storage.ErrorMethodIsNotImplemented, err)
}
