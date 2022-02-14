package serializers

import (
	"errors"

	"github.com/Fe4p3b/url-shortener/internal/serializers/json"
)

type Serializer interface {
	Encode(interface{}) ([]byte, error)
	Decode([]byte, interface{}) error
}

var _ Serializer = &json.JSONSerializer{}

func GetSerializer(t string) (Serializer, error) {
	if t == "json" {
		return &json.JSONSerializer{}, nil
	}
	return nil, errors.New("wrong type of serializer")
}
