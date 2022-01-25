package serializers

import (
	"errors"

	"github.com/Fe4p3b/url-shortener/internal/repositories"
	"github.com/Fe4p3b/url-shortener/internal/serializers/json"
	"github.com/Fe4p3b/url-shortener/internal/serializers/model"
)

type Serializer interface {
	Encode(s *model.ShortURL) ([]byte, error)
	Decode([]byte) (*model.URL, error)
	DecodeURLBatch([]byte) ([]repositories.URL, error)
	EncodeURLBatch([]repositories.URL) ([]byte, error)
}

var _ Serializer = &json.JSONSerializer{}

func GetSerializer(t string) (Serializer, error) {
	if t == "json" {
		return &json.JSONSerializer{}, nil
	}
	return nil, errors.New("wrong type of serializer")
}
