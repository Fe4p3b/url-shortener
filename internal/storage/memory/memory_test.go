package memory

import (
	"testing"

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

			assert.Equal(t, tt.want, got)
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
		url   string
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
				url:   "google.com",
				short: "qwerty",
			},
			wantErr: nil,
		},
		{
			name: "test case #2",
			args: args{
				url:   "yahoo.com",
				short: "asdf",
			},
			wantErr: storage.ErrorDuplicateShortlink,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.Save(&tt.args.short, tt.args.url)

			if err != nil && err != tt.wantErr {
				t.Errorf("Find() error = %v, wantErr = %v", err, tt.wantErr)
			}

			if _, ok := s.S[tt.args.short]; !ok && tt.wantErr == nil {
				t.Errorf("Find() value %v not found", err)
			}

		})
	}
}
