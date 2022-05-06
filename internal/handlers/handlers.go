// Package handlers provides handlers for http endpoints.
package handlers

import (
	"errors"

	"github.com/Fe4p3b/url-shortener/internal/app/shortener"
	"github.com/Fe4p3b/url-shortener/internal/models"
	"github.com/Fe4p3b/url-shortener/internal/repositories"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

var (
	ErrorURLIsGone                   = errors.New("URL is gone")
	ErrorUniqueURLViolation          = errors.New("URL already exists")
	ErrorNoContent                   = errors.New("no content")
	_                       Handlers = &handler{}
)

type Handlers interface {
	GetURL(string) (*repositories.URL, error)
	PostURL(u string, user string) (string, error)
	GetUserURLs(user string) ([]repositories.URL, error)
	DeleteUserURLs(user string, URLs []string)
	ShortenBatch(user string, batch *[]repositories.URL) ([]repositories.URL, error)
	Ping() error
	GetStats() (*models.Stats, error)
}

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

// GetURL redirects to original URL by short URL.
func (h *handler) GetURL(shortURL string) (*repositories.URL, error) {
	url, err := h.s.Find(shortURL)
	if err != nil {
		return nil, err
	}

	if url.IsDeleted {
		return nil, ErrorURLIsGone
	}

	return url, nil
}

// PostURL creates short URL by original URL.
func (h *handler) PostURL(u string, user string) (string, error) {
	sURL, err := h.s.Store(&models.URL{URL: u, UserID: user})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return "", ErrorUniqueURLViolation
		}
		return "", err
	}

	return sURL, nil
}

// GetUserURLs shows user URLs, that he created, in json.
func (h *handler) GetUserURLs(user string) ([]repositories.URL, error) {
	URLs, err := h.s.GetUserURLs(user)
	if err != nil {
		return nil, err
	}

	if len(URLs) == 0 {
		return nil, ErrorNoContent
	}

	return URLs, nil
}

// DeleteUserURLs deletes user URLs by short URL.
func (h *handler) DeleteUserURLs(user string, URLs []string) {
	h.s.DeleteURLs(user, URLs)
}

// ShortenBatch creates short URLs for batch of original URLs in json.
func (h *handler) ShortenBatch(user string, batch *[]repositories.URL) ([]repositories.URL, error) {
	sURLBatch, err := h.s.StoreBatch(user, *batch)
	if err != nil {
		return nil, err
	}

	return sURLBatch, nil
}

// Ping checks whether database connetion is up.
func (h *handler) Ping() error {
	if err := h.s.Ping(); err != nil {
		return err
	}
	return nil
}

func (h *handler) GetStats() (*models.Stats, error) {
	stats, err := h.s.GetStats()
	if err != nil {
		return nil, err
	}

	return stats, nil
}
