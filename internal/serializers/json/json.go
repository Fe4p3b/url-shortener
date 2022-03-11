// Package json provides json serializer
// for serializer package.
package json

import (
	"encoding/json"
)

// JSONSerializer encodes or decodes json data.
type JSONSerializer struct{}

// Encode implements serializers.Serializer Encode method.
func (j *JSONSerializer) Encode(v interface{}) ([]byte, error) {
	d, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return d, nil
}

// Decode implements serializers.Serializer Decode method.
func (j *JSONSerializer) Decode(b []byte, v interface{}) error {
	err := json.Unmarshal(b, v)
	if err != nil {
		return err
	}

	return nil
}
