package memory

import (
	"sync"

	"github.com/Fe4p3b/url-shortener/internal/repositories"
	"github.com/Fe4p3b/url-shortener/internal/serializers/model"
	"github.com/Fe4p3b/url-shortener/internal/storage"
)

var _ repositories.ShortenerRepository = &Memory{}

type Memory struct {
	sync.RWMutex
	S map[string]string
}

func NewMemory(s map[string]string) *Memory {
	return &Memory{
		S: s,
	}
}

func (m *Memory) Find(url string) (s string, err error) {
	v, ok := m.S[url]
	if !ok {
		return "", storage.ErrorNoLinkFound
	}
	return v, nil
}

// func (m *Memory) Save(uuid *string, url string) error {
func (m *Memory) Save(url *model.URL) error {
	if _, ok := m.S[url.ShortURL]; ok {
		return storage.ErrorDuplicateShortlink
	}

	m.Lock()
	m.S[url.ShortURL] = url.URL
	m.Unlock()
	return nil
}

func (m *Memory) GetUserURLs(user string) ([]repositories.URL, error) {
	return nil, storage.ErrorMethodIsNotImplemented
}

func (m *Memory) Ping() error {
	return storage.ErrorMethodIsNotImplemented
}

func (m *Memory) AddURLBuffer(repositories.URL) error {
	return storage.ErrorMethodIsNotImplemented
}

func (m *Memory) Flush() error {
	return storage.ErrorMethodIsNotImplemented
}
