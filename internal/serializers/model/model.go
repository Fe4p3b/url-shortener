package model

type URL struct {
	URL      string `json:"url"`
	UserId   string `json:"user_id,omitempty"`
	ShortURL string `json:"short_url"`
}

type ShortURL struct {
	ShortURL string `json:"result"`
}
