package json

import (
	"encoding/json"
	"errors"

	"github.com/Fe4p3b/url-shortener/internal/serializers/model"
)

var ErrorEmptyURL error = errors.New("url is not set")

type JSONSerializer struct{}

func (j *JSONSerializer) Encode(s *model.ShortURL) ([]byte, error) {
	d, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (j *JSONSerializer) Decode(b []byte) (*model.URL, error) {
	url := &model.URL{}
	err := json.Unmarshal(b, url)
	if err != nil {
		return nil, err
	}

	if url.URL == "" {
		return nil, ErrorEmptyURL
	}
	return url, nil
}
