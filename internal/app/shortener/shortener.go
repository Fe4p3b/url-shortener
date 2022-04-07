// Package shortener provides business logic for
// creation of shortened URLs.
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

// ShortenerService represents service for creation of short URls.
// Other than general functionality of creating
// and storing shortened URLs, the service can find original URL
// from shortened URL, can return or delete user's URLs. It also
// has ability to Ping storage for availablity.
type ShortenerService interface {
	// Find receives shortened URL and returns pointer to repositories.URL
	// if nothing was found the error is returned.
	Find(string) (*repositories.URL, error)

	// Store receives models.URL, generates short URL and tries to save it
	// in storage, if it can't be stored or short URL can't be created
	// the error is returned.
	Store(*models.URL) (string, error)

	// StoreBatch receives user identificator and repositories.URLs,
	// generates short URLs and tries to save them in a storage,
	// if repositories.URLs can't be stored or short URLs can't be created
	// the error is returned.
	StoreBatch(string, []repositories.URL) ([]repositories.URL, error)

	// GetUserURLs returns repositories.URLs for user, by user identificator,
	// or error.
	GetUserURLs(string) ([]repositories.URL, error)

	// DeleteURLs deletes URLs for user, by user identificator, asynchronously.
	DeleteURLs(string, []string)

	// Ping tests connection for the storage, or returns error.
	Ping() error

	GetStats() (*models.Stats, error)
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

// Find implements ShortenerService Find method.
func (s *shortener) Find(url string) (*repositories.URL, error) {
	return s.r.Find(url)
}

// Store implements ShortenerService Store method.
// The method generates short URL using "github.com/teris-io/shortid"
// package.
// If URL can't be saved, due to already being stored in the storage
// and postgres is used as a storage, it returns already existing URL.
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

// GetUserURLs implements ShortenerService GetUserURLs method.
func (s *shortener) GetUserURLs(user string) ([]repositories.URL, error) {
	return s.r.GetUserURLs(user, s.BaseURL)
}

// Ping implements ShortenerService Ping method.
func (s *shortener) Ping() error {
	return s.r.Ping()
}

// StoreBatch implements ShortenerService StoreBatch method.
// To optimize performance the method populates buffer of a storage,
// when buffer capacity is reached it saves all the URLs in buffer.
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

// DeleteURLs implements ShortenerService DeleteURLs method.
// To optimize performance goroutines are used. Each URL is
// added to storage channel, that serves as a buffer to
// delete URLs.
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

func (s *shortener) GetStats() (*models.Stats, error) {
	return s.r.GetStats()
}

var _ ShortenerService = &shortener{}
