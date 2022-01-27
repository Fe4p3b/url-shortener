package handlers

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/Fe4p3b/url-shortener/internal/app/shortener"
	"github.com/Fe4p3b/url-shortener/internal/middleware"
	"github.com/Fe4p3b/url-shortener/internal/serializers"
	"github.com/Fe4p3b/url-shortener/internal/serializers/json"
	"github.com/Fe4p3b/url-shortener/internal/serializers/model"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

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

func (h *handler) SetupRouting() {
	h.Router.Get("/{url}", h.GetURL)
	h.Router.Post("/", h.PostURL)
	h.Router.Post("/api/shorten", h.JSONPost)

	h.Router.Post("/api/shorten/batch", h.ShortenBatch)
	h.Router.Get("/ping", h.PingPG)

	h.Router.Get("/user/urls", h.GetUserURLs)
}

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

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *handler) PostURL(w http.ResponseWriter, r *http.Request) {
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

	user, ok := r.Context().Value(middleware.Key).(string)
	if !ok {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	sURL, err := h.s.Store(&model.URL{URL: u, UserId: user})

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

func (h *handler) JSONPost(w http.ResponseWriter, r *http.Request) {
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

	url, err := s.Decode(b)
	if err != nil {
		log.Printf("err - %v", err)
		if errors.Is(err, json.ErrorEmptyURL) {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	user, ok := r.Context().Value(middleware.Key).(string)
	if !ok {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	url.UserId = user

	log.Printf("url - %v", url)
	sURL, err := h.s.Store(url)

	var pgErr *pgconn.PgError
	header := http.StatusCreated

	if err != nil {
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			header = http.StatusConflict
		} else {
			log.Printf("error - %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	jsonSURL := &model.ShortURL{ShortURL: sURL}
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

func (h *handler) GetUserURLs(w http.ResponseWriter, r *http.Request) {
	s, err := serializers.GetSerializer("json")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// token, err := r.Cookie("token")
	// if err != nil {
	// 	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	// 	return
	// }
	user, ok := r.Context().Value(middleware.Key).(string)
	if !ok {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	log.Printf("handler - %s", user)
	URLs, err := h.s.GetUserURLs(user)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	b, err := s.EncodeURLBatch(URLs)
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

func (h *handler) ShortenBatch(w http.ResponseWriter, r *http.Request) {
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

	batch, err := s.DecodeURLBatch(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sURLBatch, err := h.s.StoreBatch(batch)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	b, err = s.EncodeURLBatch(sURLBatch)
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

func (h *handler) PingPG(w http.ResponseWriter, r *http.Request) {
	if err := h.s.Ping(); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
