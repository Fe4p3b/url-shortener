package repositories

import "github.com/Fe4p3b/url-shortener/internal/serializers/model"

type ShortenerRepository interface {
	Find(string) (string, error)
	// Save(*string, string) error
	Save(*model.URL) error
	AddURLBuffer(URL) error
	GetUserURLs(string) ([]URL, error)
	Flush() error
	Ping() error
}

type URL struct {
	CorrelationID string `json:"correlation_id,omitempty"`
	URL           string `json:"original_url,omitempty"`
	ShortURL      string `json:"short_url,omitempty"`
	UserId        string `json:"-"`
}
