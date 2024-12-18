package is

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTelNumber(t *testing.T) {
	tests := map[string]bool{
		"88888888":      true,
		"0731-88888888": true,
		"0731-8888888":  false,
		"020-88888888":  true,
		"02-88888888":   false,
	}
	for str, ok := range tests {
		err := validation.Validate(str, validation.By(TelNumber()))
		assert.Equal(t, ok, err == nil, str)
	}
}
