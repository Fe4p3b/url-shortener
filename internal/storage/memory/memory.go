// Package memory implements in-memory storage for service.
package memory

import (
	"sync"

	"github.com/Fe4p3b/url-shortener/internal/models"
	"github.com/Fe4p3b/url-shortener/internal/repositories"
	"github.com/Fe4p3b/url-shortener/internal/storage"
)

var _ repositories.ShortenerRepository = &Memory{}
var _ repositories.AuthRepository = &Memory{}

// Memory is in-memory storage.
type Memory struct {
	sync.RWMutex
	S map[string]string
}

func NewMemory(s map[string]string) *Memory {
	return &Memory{
		S: s,
	}
}

// Find implements repositories.ShortenerRepository Find method.
func (m *Memory) Find(url string) (u *repositories.URL, err error) {
	v, ok := m.S[url]

	if !ok {

		return nil, storage.ErrorNoLinkFound
	}
	u = &repositories.URL{}
	u.URL = v
	return
}

// Save implements repositories.ShortenerRepository Save method.
func (m *Memory) Save(url *models.URL) error {
	if _, ok := m.S[url.ShortURL]; ok {
		return storage.ErrorDuplicateShortlink
	}

	m.Lock()
	m.S[url.ShortURL] = url.URL
	m.Unlock()
	return nil
}

// GetUserURLs implements repositories.ShortenerRepository GetUserURLs method.
func (m *Memory) GetUserURLs(user string, baseURL string) ([]repositories.URL, error) {
	return nil, storage.ErrorMethodIsNotImplemented
}

// Ping implements repositories.ShortenerRepository Ping method.
func (m *Memory) Ping() error {
	return nil
}

// AddURLBuffer implements repositories.ShortenerRepository AddURLBuffer method.
func (m *Memory) AddURLBuffer(repositories.URL) error {
	return storage.ErrorMethodIsNotImplemented
}

// Flush implements repositories.ShortenerRepository Flush method.
func (m *Memory) Flush() error {
	return storage.ErrorMethodIsNotImplemented
}

// AddURLToDelete implements repositories.ShortenerRepository AddURLToDelete method.
func (m *Memory) AddURLToDelete(u repositories.URL) {
}

// FlushToDelete implements repositories.ShortenerRepository FlushToDelete method.
func (m *Memory) FlushToDelete() error {
	return storage.ErrorMethodIsNotImplemented
}

func (m *Memory) CreateUser() (string, error) {
	return "", storage.ErrorMethodIsNotImplemented
}
func (m *Memory) VerifyUser(string) error {
	return storage.ErrorMethodIsNotImplemented
}
