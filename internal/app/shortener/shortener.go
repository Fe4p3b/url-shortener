package shortener

import (
	"github.com/Fe4p3b/url-shortener/internal/repositories"
	"github.com/teris-io/shortid"
)

type ShortenerService interface {
	Find(string) (string, error)
	Store(string) (string, error)
	StoreBatch([]repositories.URL) ([]repositories.URL, error)
	Ping() error
}

type shortener struct {
	r repositories.ShortenerRepository
}

func NewShortener(r repositories.ShortenerRepository) *shortener {
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

func (s *shortener) Ping() error {
	return s.r.Ping()
}

func (s *shortener) StoreBatch(urls []repositories.URL) (batch []repositories.URL, err error) {
	for _, v := range urls {
		uuid, err := shortid.Generate()
		if err != nil {
			return nil, err
		}
		v.ShortURL = uuid

		if err := s.r.AddURLBuffer(v); err != nil {
			return nil, err
		}

		v.URL = ""
		batch = append(batch, v)
	}
	if err := s.r.Flush(); err != nil {
		return nil, err
	}

	return
}

var _ ShortenerService = &shortener{}
