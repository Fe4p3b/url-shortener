package repositories

type ShortenerRepository interface {
	Find(string) (string, error)
	Save(string, string) error
	AddURLBuffer(URL) error
	Flush() error
	Ping() error
}

type URL struct {
	CorrelationID string `json:"correlation_id"`
	URL           string `json:"original_url,omitempty"`
	ShortURL      string `json:"short_url,omitempty"`
}
