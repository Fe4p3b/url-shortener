package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/Fe4p3b/url-shortener/internal/app/shortener"
	"github.com/Fe4p3b/url-shortener/internal/storage/memory"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func Test_httpHandler_get(t *testing.T) {
	type fields struct {
		s      shortener.ShortenerService
		method string
		url    string
	}
	type want struct {
		code     int
		response string
		err      bool
	}

	m := memory.NewMemory(map[string]string{
		"asdf": "http://yandex.ru",
		// "qwerty": "http://google.com",
	})
	s := shortener.NewShortener(m, "http://localhost:8080")

	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "test case #1",
			fields: fields{
				s:      s,
				method: http.MethodGet,
				url:    "/asdf",
			},
			want: want{
				code:     http.StatusTemporaryRedirect,
				response: "<a href=\"http://yandex.ru\">Temporary Redirect</a>.\n\n",
				err:      false,
			},
		},
		{
			name: "test case #2",
			fields: fields{
				s:      s,
				method: http.MethodGet,
				url:    "/qwerty",
			},
			want: want{
				code:     http.StatusNotFound,
				response: "Not Found\n",
				err:      true,
			},
		},
		{
			name: "test case #3",
			fields: fields{
				s:      s,
				method: http.MethodGet,
				url:    "/",
			},
			want: want{
				code:     http.StatusNotFound,
				response: "404 page not found\n",
				err:      true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(tt.fields.s)
			request := httptest.NewRequest(tt.fields.method, tt.fields.url, nil)
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Get("/{url}", h.GetURL)
			r.ServeHTTP(w, request)

			assert.Equal(t, tt.want.code, w.Code)
			assert.Equal(t, tt.want.response, w.Body.String())

		})
	}
}

func Test_httpHandler_post(t *testing.T) {
	type fields struct {
		s      shortener.ShortenerService
		method string
		url    string
		body   string
	}
	type want struct {
		code     int
		response string
		err      bool
	}

	m := memory.NewMemory(map[string]string{
		"asdf": "yandex.ru",
	})
	s := shortener.NewShortener(m, "http://localhost:8080")

	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "test case #1",
			fields: fields{
				s:      s,
				method: http.MethodPost,
				url:    "/",
				body:   "https://yandex.ru",
			},
			want: want{
				code:     http.StatusCreated,
				response: "",
				err:      false,
			},
		},
		{
			name: "test case #2",
			fields: fields{
				s:      s,
				method: http.MethodPost,
				url:    "/",
				body:   "yandex.ru",
			},
			want: want{
				code:     http.StatusCreated,
				response: "",
				err:      false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(tt.fields.s)

			f := make(url.Values)
			f.Set("url", tt.fields.body)

			request := httptest.NewRequest(tt.fields.method, tt.fields.url, strings.NewReader(f.Encode()))
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			w := httptest.NewRecorder()

			h.PostURL(w, request)

			assert.Equal(t, tt.want.code, w.Code)
			if tt.want.response != "" {
				assert.Equal(t, tt.want.response, w.Body.String())
			}
		})
	}
}

func Test_handler_JsonPost(t *testing.T) {
	type fields struct {
		s           shortener.ShortenerService
		method      string
		url         string
		body        string
		contentType string
	}
	type want struct {
		code        int
		response    string
		err         bool
		contentType string
	}

	m := memory.NewMemory(map[string]string{
		"asdf": "yandex.ru",
	})
	s := shortener.NewShortener(m, "http://localhost:8080")

	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "test case #1",
			fields: fields{
				s:           s,
				method:      http.MethodPost,
				url:         "/api/shorten",
				body:        `{"url":"https://yandex.ru"}`,
				contentType: "application/json",
			},
			want: want{
				code:        http.StatusCreated,
				response:    "",
				err:         false,
				contentType: "application/json",
			},
		},
		{
			name: "test case #2",
			fields: fields{
				s:           s,
				method:      http.MethodPost,
				url:         "/api/shorten",
				body:        `{"url":"yandex.ru"}`,
				contentType: "application/json",
			},
			want: want{
				code:        http.StatusCreated,
				response:    "",
				err:         false,
				contentType: "application/json",
			},
		},
		{
			name: "test case #3",
			fields: fields{
				s:           s,
				method:      http.MethodPost,
				url:         "/api/shorten",
				body:        `{"url":""}`,
				contentType: "application/json",
			},
			want: want{
				code:        http.StatusBadRequest,
				response:    "Bad Request\n",
				err:         true,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(tt.fields.s)

			request := httptest.NewRequest(tt.fields.method, tt.fields.url, strings.NewReader(tt.fields.body))
			request.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			h.JSONPost(w, request)

			assert.Equal(t, tt.want.code, w.Code)
			assert.Equal(t, tt.want.contentType, w.Header().Get("Content-Type"))
			if tt.want.response != "" {
				assert.Equal(t, tt.want.response, w.Body.String())
			}
		})
	}
}
