// Package json provides json serializer
// for serializer package.
package json

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONSerializer_Encode(t *testing.T) {
	tests := []struct {
		name string
		j    *JSONSerializer
		args map[string]string
		want []byte
	}{
		{
			name: "Test case #1",
			j:    &JSONSerializer{},
			args: map[string]string{"asdf": "sdff"},
			want: []byte(`{"asdf":"sdff"}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &JSONSerializer{}
			got, err := j.Encode(tt.args)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestJSONSerializer_Decode(t *testing.T) {
	tests := []struct {
		name string
		j    *JSONSerializer
		args []byte
		want map[string]string
	}{
		{
			name: "Test case #1",
			j:    &JSONSerializer{},
			args: []byte(`{"asdf":"sdff"}`),
			want: map[string]string{"asdf": "sdff"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &JSONSerializer{}
			got := make(map[string]string)
			err := j.Decode(tt.args, &got)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
