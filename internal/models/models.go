// Package models provides required structs for URL.
package models

// URL is a struct that has original URL, short URL, and
// owner of URL
type URL struct {
	// URL is original URL
	URL string `json:"url"`

	// UserID is identification of owner
	UserID string `json:"user_id,omitempty"`

	// ShortURL is short URL
	ShortURL string `json:"short_url"`
}

// ShortURL is used for json response,
// where the key should be "result"
type ShortURL struct {
	ShortURL string `json:"result"`
}
