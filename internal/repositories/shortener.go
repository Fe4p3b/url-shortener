package repositories

type ShortenerRepository interface {
	Find(string) (string, error)
	Save(string, string) error
}
