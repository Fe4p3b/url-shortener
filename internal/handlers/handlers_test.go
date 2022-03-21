package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Fe4p3b/url-shortener/internal/app/shortener"
	"github.com/Fe4p3b/url-shortener/internal/middleware"
	"github.com/Fe4p3b/url-shortener/internal/storage/memory"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func Example() {
	// PostURL request example
	resp, err := http.Post("http://localhost:8080", "text/plain", bytes.NewReader([]byte("http://google.com")))
	if err != nil {
		fmt.Println(err)
	}
	resp.Body.Close()

	// GetURL request example
	resp, err = http.Get("http://localhost:8080/xoPnl3ang")
	if err != nil {
		fmt.Println(err)
	}
	resp.Body.Close()

	// ShortenBatch request example
	resp, err = http.Post(
		"http://localhost:8080/api/shorten/batch",
		"application/json",
		bytes.NewReader([]byte(`
		[
			{
				"correlation_id": "2765399b-d5a3-420c-8de4-f3b7fb19d334",
				"original_url": "http://aptekaplus1.kz"
			},
			{
				"correlation_id": "2765f94b-d54e-420c-8de4-f3b7fb19d325",
				"original_url": "http://hltv2.org"
			}
		]`),
		),
	)
	if err != nil {
		fmt.Println(err)
	}
	resp.Body.Close()

	// GetUserURLs request example
	resp, err = http.Get("http://localhost:8080/user/urls")
	if err != nil {
		fmt.Println(err)
	}
	resp.Body.Close()

	// JSONPost request example
	resp, err = http.Post(
		"http://localhost:8080/api/shorten",
		"application/json",
		bytes.NewReader([]byte(`
		{
			"url": "http://google.com"
		}`),
		),
	)
	if err != nil {
		fmt.Println(err)
	}
	resp.Body.Close()

	// DeleteUserURLs request example
	req, err := http.NewRequest("DELETE", "http://localhost:8080/api/user/urls", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	resp.Body.Close()

	// Ping request example
	resp, err = http.Get("http://localhost:8080/ping")
	if err != nil {
		fmt.Println(err)
	}
	resp.Body.Close()
}

func Test_handler_GetURL(t *testing.T) {
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

func Test_handler_PostURL(t *testing.T) {
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

			request := httptest.NewRequest(tt.fields.method, tt.fields.url, strings.NewReader(tt.fields.body))
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			ctx := context.WithValue(context.Background(), middleware.Key, "asdfg")

			w := httptest.NewRecorder()

			h.PostURL(w, request.WithContext(ctx))

			assert.Equal(t, tt.want.code, w.Code)
			if tt.want.response != "" {
				assert.Equal(t, tt.want.response, w.Body.String())
			}
		})
	}
}

func Test_handler_JSONPost(t *testing.T) {
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
			ctx := context.WithValue(context.Background(), middleware.Key, "asdfg")

			w := httptest.NewRecorder()

			h.JSONPost(w, request.WithContext(ctx))

			assert.Equal(t, tt.want.code, w.Code)
			assert.Equal(t, tt.want.contentType, w.Header().Get("Content-Type"))
			if tt.want.response != "" {
				assert.Equal(t, tt.want.response, w.Body.String())
			}
		})
	}
}

func Test_handler_GetUserURLs(t *testing.T) {
	m := memory.NewMemory(map[string]string{
		"asdf": "yandex.ru",
	})
	s := shortener.NewShortener(m, "http://localhost:8080")

	type fields struct {
		s           shortener.ShortenerService
		method      string
		url         string
		body        string
		contentType string
		token       string
	}
	type want struct {
		code        int
		response    string
		err         bool
		contentType string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "test case #1",
			fields: fields{
				s:           s,
				method:      http.MethodGet,
				url:         "/user/urls",
				body:        ``,
				contentType: "application/json",
				token:       "asdfg",
			},
			want: want{
				code:        http.StatusInternalServerError,
				response:    "",
				err:         true,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(tt.fields.s)

			request := httptest.NewRequest(tt.fields.method, tt.fields.url, strings.NewReader(tt.fields.body))
			request.Header.Set("Content-Type", tt.fields.contentType)
			ctx := context.WithValue(context.Background(), middleware.Key, tt.fields.token)

			w := httptest.NewRecorder()

			h.GetUserURLs(w, request.WithContext(ctx))

			assert.Equal(t, tt.want.code, w.Code)
			assert.Equal(t, tt.want.contentType, w.Header().Get("Content-Type"))
			if tt.want.response != "" {
				assert.Equal(t, tt.want.response, w.Body.String())
			}
		})
	}
}

func Test_handler_DeleteUserURLs(t *testing.T) {
	m := memory.NewMemory(map[string]string{
		"asdf": "yandex.ru",
	})
	s := shortener.NewShortener(m, "http://localhost:8080")

	type fields struct {
		s           shortener.ShortenerService
		method      string
		url         string
		body        string
		contentType string
		token       string
	}
	type want struct {
		code        int
		response    string
		err         bool
		contentType string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "test case #1",
			fields: fields{
				s:           s,
				method:      http.MethodDelete,
				url:         "/api/user/urls",
				body:        ``,
				contentType: "application/json",
				token:       "asdfg",
			},
			want: want{
				code:        http.StatusInternalServerError,
				response:    "",
				err:         true,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(tt.fields.s)

			request := httptest.NewRequest(tt.fields.method, tt.fields.url, strings.NewReader(tt.fields.body))
			request.Header.Set("Content-Type", tt.fields.contentType)
			ctx := context.WithValue(context.Background(), middleware.Key, tt.fields.token)

			w := httptest.NewRecorder()

			h.DeleteUserURLs(w, request.WithContext(ctx))

			assert.Equal(t, tt.want.code, w.Code)
			assert.Equal(t, tt.want.contentType, w.Header().Get("Content-Type"))
			if tt.want.response != "" {
				assert.Equal(t, tt.want.response, w.Body.String())
			}
		})
	}
}

func Test_handler_ShortenBatch(t *testing.T) {
	m := memory.NewMemory(map[string]string{
		"asdf": "yandex.ru",
	})
	s := shortener.NewShortener(m, "http://localhost:8080")

	type fields struct {
		s           shortener.ShortenerService
		method      string
		url         string
		body        string
		contentType string
		token       string
	}
	type want struct {
		code        int
		response    string
		err         bool
		contentType string
	}
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
				url:         "/api/shorten/batch",
				body:        ``,
				contentType: "application/json",
				token:       "asdfg",
			},
			want: want{
				code:        http.StatusBadRequest,
				response:    "",
				err:         true,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(tt.fields.s)

			request := httptest.NewRequest(tt.fields.method, tt.fields.url, strings.NewReader(tt.fields.body))
			request.Header.Set("Content-Type", tt.fields.contentType)
			ctx := context.WithValue(context.Background(), middleware.Key, tt.fields.token)

			w := httptest.NewRecorder()

			h.ShortenBatch(w, request.WithContext(ctx))

			assert.Equal(t, tt.want.code, w.Code)
			assert.Equal(t, tt.want.contentType, w.Header().Get("Content-Type"))
			if tt.want.response != "" {
				assert.Equal(t, tt.want.response, w.Body.String())
			}
		})
	}
}

func Test_handler_Ping(t *testing.T) {
	m := memory.NewMemory(map[string]string{
		"asdf": "yandex.ru",
	})
	s := shortener.NewShortener(m, "http://localhost:8080")

	type fields struct {
		s           shortener.ShortenerService
		method      string
		url         string
		body        string
		contentType string
		token       string
	}
	type want struct {
		code        int
		response    string
		err         bool
		contentType string
	}
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
				url:         "/ping",
				body:        ``,
				contentType: "application/json",
				token:       "asdfg",
			},
			want: want{
				code:        http.StatusOK,
				response:    "",
				err:         false,
				contentType: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(tt.fields.s)

			request := httptest.NewRequest(tt.fields.method, tt.fields.url, strings.NewReader(tt.fields.body))
			request.Header.Set("Content-Type", tt.fields.contentType)
			ctx := context.WithValue(context.Background(), middleware.Key, tt.fields.token)

			w := httptest.NewRecorder()

			h.Ping(w, request.WithContext(ctx))

			assert.Equal(t, tt.want.code, w.Code)
			assert.Equal(t, tt.want.contentType, w.Header().Get("Content-Type"))
			if tt.want.response != "" {
				assert.Equal(t, tt.want.response, w.Body.String())
			}
		})
	}
}
