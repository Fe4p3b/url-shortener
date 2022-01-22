package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/Fe4p3b/url-shortener/internal/app/shortener"
	"github.com/Fe4p3b/url-shortener/internal/storage/memory"
	echo "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func Test_httpHandler_get(t *testing.T) {
	type fields struct {
		s      shortener.ShortenerService
		method string
		params string
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
	s := shortener.NewShortener(m)

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
				params: "asdf",
			},
			want: want{
				code:     http.StatusTemporaryRedirect,
				response: "",
				err:      false,
			},
		},
		{
			name: "test case #2",
			fields: fields{
				s:      s,
				method: http.MethodGet,
				params: "qwerty",
			},
			want: want{
				code:     http.StatusNotFound,
				response: "Not Found",
				err:      true,
			},
		},
		{
			name: "test case #3",
			fields: fields{
				s:      s,
				method: http.MethodGet,
				params: "",
			},
			want: want{
				code:     http.StatusNotFound,
				response: "The query parameter is missing",
				err:      true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(tt.fields.s, "localhost:8080")
			request := httptest.NewRequest(tt.fields.method, "/", nil)
			w := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(request, w)

			c.SetPath("/:url")
			c.SetParamNames("url")
			c.SetParamValues(tt.fields.params)

			err := h.GetURL(c)

			if tt.want.err {
				assert.Error(t, err)
				assert.Equal(t, tt.want.code, err.(*echo.HTTPError).Code)
				assert.Equal(t, tt.want.response, err.(*echo.HTTPError).Message)
				return
			}

			assert.NoError(t, err)
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
	s := shortener.NewShortener(m)

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
			h := NewHandler(tt.fields.s, "localhost:8080")

			f := make(url.Values)
			f.Set("url", tt.fields.body)

			request := httptest.NewRequest(tt.fields.method, tt.fields.url, strings.NewReader(f.Encode()))
			request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			w := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(request, w)

			err := h.PostURL(c)
			if tt.want.err {
				assert.Error(t, err)
				assert.Equal(t, tt.want.code, err.(*echo.HTTPError).Code)
				assert.Equal(t, tt.want.response, err.(*echo.HTTPError).Message)
				return
			}

			assert.NoError(t, err)
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
	s := shortener.NewShortener(m)

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
				contentType: echo.MIMEApplicationJSON,
			},
			want: want{
				code:        http.StatusCreated,
				response:    "",
				err:         false,
				contentType: echo.MIMEApplicationJSON,
			},
		},
		{
			name: "test case #2",
			fields: fields{
				s:           s,
				method:      http.MethodPost,
				url:         "/api/shorten",
				body:        `{"url":"yandex.ru"}`,
				contentType: echo.MIMEApplicationJSON,
			},
			want: want{
				code:        http.StatusCreated,
				response:    "",
				err:         false,
				contentType: echo.MIMEApplicationJSON,
			},
		},
		{
			name: "test case #3",
			fields: fields{
				s:           s,
				method:      http.MethodPost,
				url:         "/api/shorten",
				body:        `{"url":""}`,
				contentType: echo.MIMEApplicationJSON,
			},
			want: want{
				code:        http.StatusBadRequest,
				response:    "Bad Request",
				err:         true,
				contentType: echo.MIMEApplicationJSON,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(tt.fields.s, "localhost:8080")

			request := httptest.NewRequest(tt.fields.method, tt.fields.url, strings.NewReader(tt.fields.body))
			request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			w := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(request, w)

			err := h.JSONPost(c)
			if tt.want.err {
				assert.Error(t, err)
				assert.Equal(t, tt.want.code, err.(*echo.HTTPError).Code)
				assert.Equal(t, tt.want.response, err.(*echo.HTTPError).Message)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want.code, w.Code)
			assert.Equal(t, tt.want.contentType, w.Header().Get(echo.HeaderContentType))
			if tt.want.response != "" {
				assert.Equal(t, tt.want.response, w.Body.String())
			}
		})
	}
}
