package serializers

import (
	"errors"

	"github.com/Fe4p3b/url-shortener/internal/models"
	"github.com/Fe4p3b/url-shortener/internal/repositories"
	"github.com/Fe4p3b/url-shortener/internal/serializers/json"
)

type Serializer interface {
	Encode(s *models.ShortURL) ([]byte, error)
	Decode([]byte) (*models.URL, error)
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
