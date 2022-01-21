package memory

import (
	"sync"

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

func (m *Memory) Find(url string) (s string, err error) {
	v, ok := m.S[url]
	if !ok {
		return "", storage.ErrorNoLinkFound
	}
	return v, nil
}

func (m *Memory) Save(uuid string, url string) error {
	if _, ok := m.S[uuid]; ok {
		return storage.ErrorDuplicateShortlink
	}

	m.Lock()
	m.S[uuid] = url
	m.Unlock()
	return nil
}
