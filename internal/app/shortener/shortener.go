package shortener

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/Fe4p3b/url-shortener/internal/models"
	"github.com/Fe4p3b/url-shortener/internal/repositories"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/teris-io/shortid"

	"golang.org/x/sync/errgroup"
)

type ShortenerService interface {
	Find(string) (*repositories.URL, error)
	Store(*models.URL) (string, error)
	StoreBatch(string, []repositories.URL) ([]repositories.URL, error)
	GetUserURLs(string) ([]repositories.URL, error)
	DeleteURLs(string, []string)
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

func (s *shortener) Find(url string) (*repositories.URL, error) {
	return s.r.Find(url)
}

func (s *shortener) Store(url *models.URL) (string, error) {
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
	return s.r.GetUserURLs(user, s.BaseURL)
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
		v.UserID = user

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

func (s *shortener) DeleteURLs(user string, URLs []string) {
	go func() {
		g, _ := errgroup.WithContext(context.Background())

		g.Go(func() error {
			if err := s.r.FlushToDelete(); err != nil {
				return err
			}

			return nil
		})

		wg := &sync.WaitGroup{}
		for _, url := range URLs {
			wg.Add(1)
			go func(url string) {
				defer wg.Done()
				s.r.AddURLToDelete(repositories.URL{ShortURL: url, UserID: user})
			}(url)
		}
		wg.Wait()

		if err := g.Wait(); err != nil {
			log.Printf("error deleting urls for user: %v", err)
		}
	}()

}

var _ ShortenerService = &shortener{}
