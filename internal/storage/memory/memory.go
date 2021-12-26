package memory

import (
	"errors"
	"sync"

	"github.com/Fe4p3b/url-shortener/internal/app/shortener"
)

var _ shortener.ShortenerRepository = &memory{}

type memory struct {
	sync.RWMutex
	s map[string]string
}

func New() *memory {
	return &memory{
		s: make(map[string]string),
	}
}

func (m *memory) Find(url string) (s string, err error) {
	v, ok := m.s[url]
	if !ok {
		return "", errors.New("No such link")
	}
	return v, nil
}

func (m *memory) Save(uuid string, url string) error {
	if _, ok := m.s[uuid]; ok {
		return errors.New("Couldnt save, duplicate shortlink")
	}

	m.Lock()
	m.s[uuid] = url
	m.Unlock()
	return nil
}
