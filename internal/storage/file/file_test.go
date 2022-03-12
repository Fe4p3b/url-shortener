package file

import (
	"os"
	"testing"

	"github.com/Fe4p3b/url-shortener/internal/repositories"
	"github.com/Fe4p3b/url-shortener/internal/storage"
	"github.com/Fe4p3b/url-shortener/internal/storage/memory"
	"github.com/stretchr/testify/assert"
)

func Test_file_Find(t *testing.T) {
	s := &memory.Memory{
		S: map[string]string{
			"asdf": "yandex.ru",
		},
	}
	f := file{
		m: s,
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
			got, err := f.Find(tt.value)

			if err != nil && tt.wantErr != err {
				assert.Equal(t, tt.want, got)
			}

			if got != nil {
				assert.Equal(t, tt.want, got.URL)
			}
		})
	}
}

func TestMemory_Ping(t *testing.T) {
	f, err := NewFile("test")
	if err != nil {
		t.Error(err)
	}
	err = f.Ping()
	assert.NoError(t, err)

	if err := os.Remove("test"); err != nil {
		t.Error(err)
	}
}

func TestMemory_Close(t *testing.T) {
	f, err := NewFile("test")
	if err != nil {
		t.Error(err)
	}
	err = f.Close()
	assert.NoError(t, err)

	if err := os.Remove("test"); err != nil {
		t.Error(err)
	}
}

func TestMemory_GetUserURLs(t *testing.T) {
	s := &memory.Memory{
		S: map[string]string{
			"asdf": "yandex.ru",
		},
	}
	f := file{
		m: s,
	}
	_, err := f.GetUserURLs("", "")
	assert.Error(t, storage.ErrorMethodIsNotImplemented, err)
}

func TestMemory_AddURLBuffer(t *testing.T) {
	s := &memory.Memory{
		S: map[string]string{
			"asdf": "yandex.ru",
		},
	}
	f := file{
		m: s,
	}
	err := f.AddURLBuffer(repositories.URL{})
	assert.Error(t, storage.ErrorMethodIsNotImplemented, err)
}

func TestMemory_Flush(t *testing.T) {
	s := &memory.Memory{
		S: map[string]string{
			"asdf": "yandex.ru",
		},
	}
	f := file{
		m: s,
	}
	err := f.Flush()
	assert.Error(t, storage.ErrorMethodIsNotImplemented, err)
}

func TestMemory_FlushToDelete(t *testing.T) {
	s := &memory.Memory{
		S: map[string]string{
			"asdf": "yandex.ru",
		},
	}
	f := file{
		m: s,
	}
	err := f.FlushToDelete()
	assert.Error(t, storage.ErrorMethodIsNotImplemented, err)
}
