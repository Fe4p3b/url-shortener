// Package serializers provides serialization
// functionality.
package serializers

import (
	"testing"

	"github.com/Fe4p3b/url-shortener/internal/serializers/json"
	"github.com/stretchr/testify/assert"
)

func TestGetSerializer(t *testing.T) {

	tests := []struct {
		name           string
		serializerType string
		wantErr        bool
	}{
		{
			name:           "Test case #1",
			serializerType: "bson",
			wantErr:        true,
		},
		{
			name:           "Test case #2",
			serializerType: "json",
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := GetSerializer(tt.serializerType)
			if err != nil {
				assert.Error(t, ErrorSerializerType, err)
			} else {
				assert.IsType(t, &json.JSONSerializer{}, s)
			}
		})
	}
}
