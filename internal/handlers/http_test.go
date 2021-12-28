package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/Fe4p3b/url-shortener/internal/app/shortener"
	"github.com/Fe4p3b/url-shortener/internal/storage/memory"
)

func Test_httpHandler_handler(t *testing.T) {
	type fields struct {
		s      shortener.ShortenerService
		method string
		url    string
		body   url.Values
	}
	type want struct {
		code     int
		response string
	}

	m := memory.New(map[string]string{
		"asdf": "yandex.ru",
	})
	s := shortener.New(m)

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
				body: url.Values{
					"url": []string{"https://yandex.ru"},
				},
			},
			want: want{
				code:     http.StatusCreated,
				response: "",
			},
		},
		{
			name: "test case #2",
			fields: fields{
				s:      s,
				method: http.MethodPost,
				url:    "/",
				body: url.Values{
					"url": []string{"yandex.ru"},
				},
			},
			want: want{
				code:     http.StatusInternalServerError,
				response: "Invalid URL\n\n",
			},
		},
		{
			name: "test case #3",
			fields: fields{
				s:      s,
				method: http.MethodGet,
				url:    "/?url=asdf",
				body:   nil,
			},
			want: want{
				code:     http.StatusTemporaryRedirect,
				response: "<a href=\"/yandex.ru\">Temporary Redirect</a>.\n\n\n",
			},
		},
		{
			name: "test case #4",
			fields: fields{
				s:      s,
				method: http.MethodGet,
				url:    "/",
				body:   nil,
			},
			want: want{
				code:     http.StatusBadRequest,
				response: "The query parameter is missing\n\n",
			},
		},
		{
			name: "test case #5",
			fields: fields{
				s:      s,
				method: http.MethodPut,
				url:    "/",
				body:   nil,
			},
			want: want{
				code:     http.StatusMethodNotAllowed,
				response: "\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpHandler := &httpHandler{
				s: tt.fields.s,
			}
			request := httptest.NewRequest(tt.fields.method, tt.fields.url, strings.NewReader(tt.fields.body.Encode()))
			if tt.fields.body != nil {
				request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			}
			w := httptest.NewRecorder()
			h := http.HandlerFunc(httpHandler.handler)
			h.ServeHTTP(w, request)
			res := w.Result()
			if res.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d, got %d", tt.want.code, w.Code)
			}

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}

			// t.Errorf("%v", tt.want.response)
			// t.Errorf("%v", string(resBody))
			// t.Errorf("%v", reflect.DeepEqual(tt.want.response, string(resBody)))
			// t.Errorf("Expected body %v, got %v => %v", (tt.want.response), string(resBody), string(resBody) != tt.want.response)
			if string(resBody) != tt.want.response && tt.want.response != "" {
				t.Errorf("Expected body %s, got %s => %v --- %v", tt.want.response, string(resBody), string(resBody) != tt.want.response, res.Header.Get("Content-Type"))
			}
		})
	}
}
