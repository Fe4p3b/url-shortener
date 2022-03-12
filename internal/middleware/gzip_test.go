package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGZIPReaderMiddleware(t *testing.T) {
	type fields struct {
		next            http.Handler
		acceptEncoding  string
		contentEncoding string
		body            string
	}
	type want struct {
		code     int
		response string
		err      bool
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "Test case #1",
			fields: fields{
				next:            http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
				acceptEncoding:  "gzip",
				contentEncoding: "gzip",
				body:            "",
			},
			want: want{
				code:     http.StatusBadRequest,
				response: "Bad Request",
			},
		},
		{
			name: "Test case #2",
			fields: fields{
				next:            http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
				acceptEncoding:  "gzip",
				contentEncoding: "gzip",
				body:            "\x1f\x8b\b\x00\x00\x00\x00\x00\x04\xff\x01\x00\x00\xff\xff\x00\x00\x00\x00\x00\x00\x00\x00",
			},
			want: want{
				code:     http.StatusOK,
				response: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(tt.fields.body))
			request.Header.Set("Content-Encoding", tt.fields.contentEncoding)
			request.Header.Set("Accept-Encoding", tt.fields.acceptEncoding)
			w := httptest.NewRecorder()

			GZIPReaderMiddleware(tt.fields.next).ServeHTTP(w, request)
			assert.Equal(t, tt.want.code, w.Code)
			if tt.want.response != "" {
				assert.Equal(t, tt.want.response, w.Body.String())
			}
		})
	}
}

func TestGZIPWriterMiddleware(t *testing.T) {
	type fields struct {
		next            http.Handler
		acceptEncoding  string
		contentEncoding string
		body            string
	}
	type want struct {
		code     int
		response string
		err      bool
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "Test case #1",
			fields: fields{
				next:            http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
				acceptEncoding:  "gzip",
				contentEncoding: "gzip",
				body:            "asdf",
			},
			want: want{
				code:     http.StatusOK,
				response: "\x1f\x8b\b\x00\x00\x00\x00\x00\x04\xff\x01\x00\x00\xff\xff\x00\x00\x00\x00\x00\x00\x00\x00",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(tt.fields.body))
			request.Header.Set("Content-Encoding", tt.fields.contentEncoding)
			request.Header.Set("Accept-Encoding", tt.fields.acceptEncoding)
			w := httptest.NewRecorder()

			GZIPWriterMiddleware(tt.fields.next).ServeHTTP(w, request)
			assert.Equal(t, tt.want.code, w.Code)
			if tt.want.response != "" {
				assert.Equal(t, tt.want.response, w.Body.String())
			}
		})
	}
}
