package is

import (
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/stretchr/testify/assert"
)

func TestMillisecond(t *testing.T) {
	tests := map[int64]bool{
		0:             false,
		23:            false,
		732855781313:  false,
		1732855781313: true,
		1733101007000: true,
		1733101007356: true,
	}
	for str, ok := range tests {
		err := validation.Validate(str, validation.By(Millisecond()))
		assert.Equal(t, ok, err == nil, str)
	}
}
