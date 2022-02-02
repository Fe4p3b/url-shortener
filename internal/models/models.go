package models

type URL struct {
	URL      string `json:"url"`
	UserID   string `json:"user_id,omitempty"`
	ShortURL string `json:"short_url"`
}

type ShortURL struct {
	ShortURL string `json:"result"`
}
