package json

import (
	"encoding/json"
	"errors"

	"github.com/Fe4p3b/url-shortener/internal/models"
	"github.com/Fe4p3b/url-shortener/internal/repositories"
)

var ErrorEmptyURL error = errors.New("url is not set")

type JSONSerializer struct{}

func (j *JSONSerializer) Encode(s *models.ShortURL) ([]byte, error) {
	d, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (j *JSONSerializer) Decode(b []byte) (*models.URL, error) {
	url := &models.URL{}
	err := json.Unmarshal(b, url)
	if err != nil {
		return nil, err
	}

	if url.URL == "" {
		return nil, ErrorEmptyURL
	}
	return url, nil
}

func (j *JSONSerializer) DecodeURLBatch(b []byte) (batch []repositories.URL, err error) {
	err = json.Unmarshal(b, &batch)
	if err != nil {
		return nil, err
	}
	return batch, nil
}

func (j *JSONSerializer) EncodeURLBatch(batch []repositories.URL) (b []byte, err error) {
	b, err = json.Marshal(batch)
	if err != nil {
		return nil, err
	}
	return b, nil
}
