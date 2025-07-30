package temu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_isf(t *testing.T) {
	testFormats := map[string]bool{
		"Hello, %s!":       true,
		"Value: %d":        true,
		"Percent: %%":      true,
		"Complex: %+10.2f": true,
		"Float: %.2f":      true,
		"Invalid: %z":      false,
		"Incomplete: %":    false,
		"Incomplete: %%%":  false,
	}

	for format, ok := range testFormats {
		assert.Equal(t, ok, isf(format), format)
	}
}
