package file

import (
	"bufio"
	"errors"
	"io"
	"os"

	"github.com/Fe4p3b/url-shortener/internal/repositories"
	"github.com/Fe4p3b/url-shortener/internal/storage/memory"
	"gopkg.in/yaml.v2"
)

type file struct {
	file *os.File
	rw   *bufio.ReadWriter
	m    *memory.Memory
}

var _ repositories.ShortenerRepository = &file{}

func NewFile(path string) (*file, error) {
	m := memory.NewMemory(map[string]string{})
	s := &file{
		m: m,
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

func (f *file) Find(url string) (s string, err error) {
	s, err = f.m.Find(url)
	return
}

func (f *file) Save(uuid string, url string) error {
	if err := f.m.Save(uuid, url); err != nil {
		return err
	}

	data, err := yaml.Marshal(&map[string]string{uuid: url})
	if err != nil {
		return err
	}

	if _, err := f.rw.Writer.Write(data); err != nil {
		return err
	}

	return f.rw.Writer.Flush()
}

func (f *file) Close() error {
	return f.file.Close()
}

func (f *file) Ping() error {
	return nil
}
