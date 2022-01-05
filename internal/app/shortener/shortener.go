package shortener

import (
	"github.com/Fe4p3b/url-shortener/internal/repositories"
	"github.com/teris-io/shortid"
)

type ShortenerService interface {
	Find(string) (string, error)
	Store(string) (string, error)
}

type shortener struct {
	r repositories.ShortenerRepository
}

func New(r repositories.ShortenerRepository) *shortener {
	return &shortener{
		r: r,
	}
}

func (s *shortener) Find(url string) (string, error) {
	return s.r.Find(url)
}

func (s *shortener) Store(url string) (string, error) {
	uuid, err := shortid.Generate()
	if err != nil {
		return "", err
	}
	err = s.r.Save(uuid, url)
	if err != nil {
		return "", err
	}
	return uuid, nil
}
