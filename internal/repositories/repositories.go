// Package repositories provides interfaces and structs
// for storages to use and implement.
package repositories

import "github.com/Fe4p3b/url-shortener/internal/models"

// ShortenerRepository provides functionality to find,
// store and delete from storage.
type ShortenerRepository interface {
	// Find finds URL by short URL.
	Find(string) (*URL, error)

	// Save stores models.URL in a storage.
	Save(*models.URL) error

	// AddURLBuffer adds URL to add buffer, that is used
	// to optimize bulk URL addition.
	AddURLBuffer(URL) error

	// AddURLToDelete URL to delete buffer, that is used
	// to optimize URL deletion.
	AddURLToDelete(URL)

	// GetUserURLs return slice of URLs for user, with
	// certain base URL, like localhost:8080.
	GetUserURLs(string, string) ([]URL, error)

	// Flush flushes buffer, that is used for bulk URL
	// addition
	Flush() error

	// FlushToDelete flushes delete buffer, that is used for
	// optimized URL deletion.
	FlushToDelete() error

	// Ping tests connection with storage or returns error.
	Ping() error
}

// AuthRepository provides functionality to create user identificator
// or verify whether user exists in storage.
type AuthRepository interface {
	CreateUser() (string, error)
	VerifyUser(string) error
}

// URL is used to store or retrive bulk data from storage.
type URL struct {
	CorrelationID string `json:"correlation_id,omitempty"`
	URL           string `json:"original_url,omitempty"`
	ShortURL      string `json:"short_url,omitempty"`
	UserID        string `json:"-"`
	IsDeleted     bool   `json:"-"`
}
