package json

import (
	"encoding/json"
)

type JSONSerializer struct{}

func (j *JSONSerializer) Encode(v interface{}) ([]byte, error) {
	d, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (j *JSONSerializer) Decode(b []byte, v interface{}) error {
	err := json.Unmarshal(b, v)
	if err != nil {
		return err
	}

	return nil
}
