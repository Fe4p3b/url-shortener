package memory

import (
	"log"
	"sync"

	"github.com/Fe4p3b/url-shortener/internal/models"
	"github.com/Fe4p3b/url-shortener/internal/repositories"
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

func (m *Memory) Find(url string) (u *repositories.URL, err error) {
	v, ok := m.S[url]
	log.Println(v)

	if !ok {
		log.Println(ok)

		return nil, storage.ErrorNoLinkFound
	}
	u = &repositories.URL{}
	u.URL = v
	log.Println(u)
	return
}

func (m *Memory) Save(url *models.URL) error {
	if _, ok := m.S[url.ShortURL]; ok {
		return storage.ErrorDuplicateShortlink
	}

	m.Lock()
	m.S[url.ShortURL] = url.URL
	m.Unlock()
	return nil
}

func (m *Memory) GetUserURLs(user string, baseURL string) ([]repositories.URL, error) {
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

func (m *Memory) AddURLToDelete(u repositories.URL) {
}

func (m *Memory) FlushToDelete() error {
	return storage.ErrorMethodIsNotImplemented
}
