package json

import (
	"encoding/json"
	"errors"

	"github.com/Fe4p3b/url-shortener/internal/serializers/model"
)

var ErrorEmptyUrl error = errors.New("url is not set")

type JsonSerializer struct{}

func (j *JsonSerializer) Encode(s *model.SURL) ([]byte, error) {
	d, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (j *JsonSerializer) Decode(b []byte) (*model.Url, error) {
	url := &model.Url{}
	err := json.Unmarshal(b, url)
	if err != nil {
		return nil, err
	}

	if url.Url == "" {
		return nil, ErrorEmptyUrl
	}
	return url, nil
}
