package handlers

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/Fe4p3b/url-shortener/internal/app/shortener"
	"github.com/Fe4p3b/url-shortener/internal/serializers"
	"github.com/Fe4p3b/url-shortener/internal/serializers/json"
	"github.com/Fe4p3b/url-shortener/internal/serializers/model"
	"github.com/go-chi/chi/v5"
)

type handler struct {
	s       shortener.ShortenerService
	BaseURL string
	Router  *chi.Mux
}

func NewHandler(s shortener.ShortenerService, bURL string) *handler {
	return &handler{
		s:       s,
		BaseURL: bURL,
		Router:  chi.NewRouter(),
	}
}

func (h *handler) SetupRouting() {
	h.Router.Get("/{url}", h.GetURL)
	h.Router.Post("/", h.PostURL)
	h.Router.Post("/api/shorten", h.JSONPost)
}

func (h *handler) GetURL(w http.ResponseWriter, r *http.Request) {
	q := chi.URLParam(r, "url")

	if q == "" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("The query parameter is missing"))
		return
	}

	url, err := h.s.Find(q)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(http.StatusText(http.StatusNotFound)))
		return
	}

	log.Printf("query_url - %s, url - %s, status - %d", q, url, http.StatusTemporaryRedirect)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *handler) PostURL(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}

	u := r.Form.Get("url")

	_, err := url.Parse(u)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(http.StatusText(http.StatusBadRequest)))
		return
	}

	sURL, err := h.s.Store(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}

	log.Printf("query_url - %s, short-url - %s", u, sURL)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("%s/%s", h.BaseURL, sURL)))
}

func (h *handler) JSONPost(w http.ResponseWriter, r *http.Request) {
	s, err := serializers.GetSerializer("json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}

	url, err := s.Decode(b)
	if err != nil {
		if errors.Is(err, json.ErrorEmptyURL) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(http.StatusText(http.StatusBadRequest)))
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}

	sURL, err := h.s.Store(url.URL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}

	log.Printf("query_url - %s, short-url - %s", url, sURL)
	jsonSURL := &model.ShortURL{ShortURL: fmt.Sprintf("%s/%s", h.BaseURL, sURL)}
	b, err = s.Encode(jsonSURL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(b)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}

}
