// Package handlers provides handlers for http endpoints.
package handlers

import (
	"errors"
	"io"
	"net/http"
	"net/url"

	"net/http/pprof"

	"github.com/Fe4p3b/url-shortener/internal/app/shortener"
	"github.com/Fe4p3b/url-shortener/internal/middleware"
	"github.com/Fe4p3b/url-shortener/internal/models"
	"github.com/Fe4p3b/url-shortener/internal/repositories"
	"github.com/Fe4p3b/url-shortener/internal/serializers"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

// handler provides handlers for http endpoints.
type handler struct {
	s      shortener.ShortenerService
	Router *chi.Mux
}

func NewHandler(s shortener.ShortenerService) *handler {
	return &handler{
		s:      s,
		Router: chi.NewRouter(),
	}
}

// SetupAPIRouting initializes http routes for api.
func (h *handler) SetupAPIRouting() {
	h.Router.Get("/{url}", h.GetURL)
	h.Router.Post("/", h.PostURL)
	h.Router.Post("/api/shorten", h.JSONPost)

	h.Router.Post("/api/shorten/batch", h.ShortenBatch)
	h.Router.Get("/ping", h.Ping)

	h.Router.Get("/user/urls", h.GetUserURLs)
	h.Router.Delete("/api/user/urls", h.DeleteUserURLs)
}

// SetupProfiling initializes http routes for profiling.
func (h *handler) SetupProfiling() {
	h.Router.HandleFunc("/debug/pprof/", pprof.Index)
	h.Router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	h.Router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	h.Router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	h.Router.HandleFunc("/debug/pprof/trace", pprof.Trace)

	h.Router.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	h.Router.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	h.Router.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))
	h.Router.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	h.Router.Handle("/debug/pprof/block", pprof.Handler("block"))
	h.Router.Handle("/debug/pprof/allocs", pprof.Handler("allocs"))
}

// GetURL redirects to original URL by short URL.
func (h *handler) GetURL(w http.ResponseWriter, r *http.Request) {
	q := chi.URLParam(r, "url")

	if q == "" {
		http.Error(w, "The query parameter is missing", http.StatusNotFound)
		return
	}

	url, err := h.s.Find(q)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if url.IsDeleted {
		http.Error(w, http.StatusText(http.StatusGone), http.StatusGone)
		return
	}

	http.Redirect(w, r, url.URL, http.StatusTemporaryRedirect)
}

// PostURL creates short URL by original URL.
func (h *handler) PostURL(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.Key).(string)
	if !ok {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	u := string(b)

	_, err = url.Parse(u)
	if err != nil || len(u) == 0 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	sURL, err := h.s.Store(&models.URL{URL: u, UserID: user})

	var pgErr *pgconn.PgError
	header := http.StatusCreated

	if err != nil {
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			header = http.StatusConflict
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(header)
	_, err = w.Write([]byte(sURL))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

// JSONPost creates short URL by original URL in json.
func (h *handler) JSONPost(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.Key).(string)
	if !ok {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	s, err := serializers.GetSerializer("json")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	url := &models.URL{}
	if err = s.Decode(b, url); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if url.URL == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	url.UserID = user

	sURL, err := h.s.Store(url)

	var pgErr *pgconn.PgError
	header := http.StatusCreated

	if err != nil {
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			header = http.StatusConflict
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	jsonSURL := &models.ShortURL{ShortURL: sURL}
	b, err = s.Encode(jsonSURL)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(header)
	_, err = w.Write(b)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

// GetUserURLs shows user URLs, that he created, in json.
func (h *handler) GetUserURLs(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.Key).(string)
	if !ok {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	s, err := serializers.GetSerializer("json")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	URLs, err := h.s.GetUserURLs(user)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if len(URLs) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	b, err := s.Encode(URLs)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(b)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

// DeleteUserURLs deletes user URLs by short URL.
func (h *handler) DeleteUserURLs(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.Key).(string)
	if !ok {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	s, err := serializers.GetSerializer("json")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	URLs := make([]string, 0)

	if err := s.Decode(b, &URLs); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	h.s.DeleteURLs(user, URLs)

	w.WriteHeader(http.StatusAccepted)
}

// ShortenBatch creates short URLs for batch of original URLs in json.
func (h *handler) ShortenBatch(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.Key).(string)
	if !ok {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	s, err := serializers.GetSerializer("json")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	batch := &[]repositories.URL{}
	if err = s.Decode(b, batch); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sURLBatch, err := h.s.StoreBatch(user, *batch)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	b, err = s.Encode(sURLBatch)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(b)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

// Ping checks whether database connetion is up.
func (h *handler) Ping(w http.ResponseWriter, r *http.Request) {
	if err := h.s.Ping(); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
