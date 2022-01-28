package shortener

import (
	"errors"
	"fmt"

	"github.com/Fe4p3b/url-shortener/internal/repositories"
	"github.com/Fe4p3b/url-shortener/internal/serializers/model"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/teris-io/shortid"
)

type ShortenerService interface {
	Find(string) (string, error)
	Store(*model.URL) (string, error)
	StoreBatch(string, []repositories.URL) ([]repositories.URL, error)
	GetUserURLs(string) ([]repositories.URL, error)
	Ping() error
}

type shortener struct {
	r       repositories.ShortenerRepository
	BaseURL string
}

func NewShortener(r repositories.ShortenerRepository, u string) *shortener {
	return &shortener{
		r:       r,
		BaseURL: u,
	}
}

func (s *shortener) Find(url string) (string, error) {
	return s.r.Find(url)
}

func (s *shortener) Store(url *model.URL) (string, error) {
	uuid, err := shortid.Generate()
	if err != nil {
		return "", err
	}

	url.ShortURL = uuid
	err = s.r.Save(url)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return fmt.Sprintf("%s/%s", s.BaseURL, url.ShortURL), err
		}
		return "", err
	}

	return fmt.Sprintf("%s/%s", s.BaseURL, url.ShortURL), nil
}

func (s *shortener) GetUserURLs(user string) ([]repositories.URL, error) {
	return s.r.GetUserURLs(user)
}

func (s *shortener) Ping() error {
	return s.r.Ping()
}

func (s *shortener) StoreBatch(user string, urls []repositories.URL) (batch []repositories.URL, err error) {
	for _, v := range urls {
		uuid, err := shortid.Generate()
		if err != nil {
			return nil, err
		}
		v.ShortURL = uuid
		v.UserId = user

		if err := s.r.AddURLBuffer(v); err != nil {
			return nil, err
		}

		v.URL = ""
		v.ShortURL = fmt.Sprintf("%s/%s", s.BaseURL, uuid)
		batch = append(batch, v)
	}
	if err := s.r.Flush(); err != nil {
		return nil, err
	}

	return
}

var _ ShortenerService = &shortener{}
