package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTruncateWeightValue(t *testing.T) {
	tests := []struct {
		name  string
		value int64
		want  int64
	}{
		{"1", 0, 0},
		{"2", 1, 1000},
		{"3", 999, 1000},
		{"4", 1000, 1000},
		{"5", 1001, 2000},
		{"6", 1100, 2000},
		{"7", 1999, 2000},
		{"8", 2000, 2000},
		{"9", 200001, 201000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, TruncateWeightValue(tt.value), "TruncateWeightValue(%v)", tt.value)
		})
	}
}
