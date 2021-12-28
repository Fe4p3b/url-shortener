package memory

import (
	"errors"
	"sync"

	"github.com/Fe4p3b/url-shortener/internal/repositories"
)

var noLinkFoundError = errors.New("No such link")
var duplicateShortlinkError = errors.New("No such link")
var _ repositories.ShortenerRepository = &memory{}

type memory struct {
	sync.RWMutex
	S map[string]string
}

func New(s map[string]string) *memory {
	return &memory{
		S: s,
	}
}

func (m *memory) Find(url string) (s string, err error) {
	v, ok := m.S[url]
	if !ok {
		return "", noLinkFoundError
	}
	return v, nil
}

func (m *memory) Save(uuid string, url string) error {
	if _, ok := m.S[uuid]; ok {
		return duplicateShortlinkError
	}

	m.Lock()
	m.S[uuid] = url
	m.Unlock()
	return nil
}
