package serializers

import (
	"errors"

	"github.com/Fe4p3b/url-shortener/internal/serializers/json"
	"github.com/Fe4p3b/url-shortener/internal/serializers/model"
)

type FactorySerializer interface {
	Encode(s *model.SURL) ([]byte, error)
	Decode([]byte) (*model.Url, error)
}

var _ FactorySerializer = &json.JsonSerializer{}

func GetSerializer(t string) (FactorySerializer, error) {
	if t == "json" {
		return &json.JsonSerializer{}, nil
	}
	return nil, errors.New("wrong type of serializer")
}
