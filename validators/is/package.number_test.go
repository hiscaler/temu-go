package is

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPackageNumber(t *testing.T) {
	tests := map[string]bool{
		"pc123":            false,
		" pc123":           false,
		"pc1234567890":     false,
		"PC2411151434535":  true,
		"pc2411151434535":  true,
		"pC2411151434535":  true,
		"Pc2411151434535":  true,
		"PC24111514345351": false,
		"PC2411151434535 ": false,
		" PC2411151434535": false,
	}
	for str, ok := range tests {
		err := validation.Validate(str, validation.By(PackageNumber()))
		assert.Equal(t, ok, err == nil, str)
	}
}
