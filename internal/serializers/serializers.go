// Package serializers provides serialization
// functionality.
package serializers

import (
	"errors"

	"github.com/Fe4p3b/url-shortener/internal/serializers/json"
)

var (
	ErrorSerializerType = errors.New("wrong type of serializer")
)

// Serializer encodes or decodes data.
type Serializer interface {
	// Encode encodes interface to slice of bytes.
	Encode(interface{}) ([]byte, error)

	// Decode decodes slice of bytes into the interface.
	Decode([]byte, interface{}) error
}

var _ Serializer = &json.JSONSerializer{}

func GetSerializer(t string) (Serializer, error) {
	if t == "json" {
		return &json.JSONSerializer{}, nil
	}
	return nil, ErrorSerializerType
}
