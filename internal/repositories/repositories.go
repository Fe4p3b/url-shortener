package repositories

import "github.com/Fe4p3b/url-shortener/internal/models"

type ShortenerRepository interface {
	Find(string) (string, error)
	Save(*models.URL) error
	AddURLBuffer(URL) error
	GetUserURLs(string) ([]URL, error)
	Flush() error
	Ping() error
}

type AuthRepository interface {
	CreateUser() (string, error)
	VerifyUser(string) error
}

type URL struct {
	CorrelationID string `json:"correlation_id,omitempty"`
	URL           string `json:"original_url,omitempty"`
	ShortURL      string `json:"short_url,omitempty"`
	UserID        string `json:"-"`
}
