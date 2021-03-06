// Package file implements file storage for service.
package file

import (
	"bufio"
	"errors"
	"io"
	"os"

	"github.com/Fe4p3b/url-shortener/internal/models"
	"github.com/Fe4p3b/url-shortener/internal/repositories"
	"github.com/Fe4p3b/url-shortener/internal/storage"
	"github.com/Fe4p3b/url-shortener/internal/storage/memory"
	"gopkg.in/yaml.v2"
)

// file implements file storage.
type file struct {
	path string
	file *os.File
	rw   *bufio.ReadWriter
	m    *memory.Memory
}

var _ repositories.ShortenerRepository = &file{}

func NewFile(path string) (*file, error) {
	m := memory.NewMemory(map[string]string{})
	s := &file{
		path: path,
		m:    m,
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0777)
	if err == nil {
		s.file = f
		s.rw = bufio.NewReadWriter(bufio.NewReader(f), bufio.NewWriter(f))

		data, err := io.ReadAll(s.rw.Reader)
		if err != nil {
			return nil, err
		}

		if err = yaml.Unmarshal(data, &s.m.S); err != nil {
			return nil, err
		}

		return s, nil
	}

	if !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	f, err = os.Create(path)
	if err != nil {
		return nil, err
	}

	s.file = f
	s.rw = bufio.NewReadWriter(bufio.NewReader(f), bufio.NewWriter(f))
	return s, nil
}

// Find implements repositories.ShortenerRepository Find method.
func (f *file) Find(url string) (u *repositories.URL, err error) {
	u, err = f.m.Find(url)
	return
}

// Save implements repositories.ShortenerRepository Save method.
func (f *file) Save(url *models.URL) error {
	if err := f.m.Save(url); err != nil {
		return err
	}

	data, err := yaml.Marshal(&map[string]string{url.ShortURL: url.URL})
	if err != nil {
		return err
	}

	if _, err := f.rw.Writer.Write(data); err != nil {
		return err
	}

	return f.rw.Writer.Flush()
}

// Close closes file.
func (f *file) Close() error {
	return f.file.Close()
}

// GetUserURLs implements repositories.ShortenerRepository GetUserURLs method.
func (f *file) GetUserURLs(user string, baseURL string) ([]repositories.URL, error) {
	return nil, storage.ErrorMethodIsNotImplemented
}

// Ping implements repositories.ShortenerRepository Ping method.
func (f *file) Ping() error {
	if _, err := os.Stat(f.path); os.IsNotExist(err) {
		return err
	}
	return nil
}

// AddURLBuffer implements repositories.ShortenerRepository AddURLBuffer method.
func (f *file) AddURLBuffer(repositories.URL) error {
	return storage.ErrorMethodIsNotImplemented
}

// Flush implements repositories.ShortenerRepository Flush method.
func (f *file) Flush() error {
	return storage.ErrorMethodIsNotImplemented
}

// AddURLToDelete implements repositories.ShortenerRepository AddURLToDelete method.
func (f *file) AddURLToDelete(u repositories.URL) {
}

// FlushToDelete implements repositories.ShortenerRepository FlushToDelete method.
func (f *file) FlushToDelete() error {
	return storage.ErrorMethodIsNotImplemented
}

func (f *file) GetStats() (*models.Stats, error) {
	return nil, storage.ErrorMethodIsNotImplemented
}
