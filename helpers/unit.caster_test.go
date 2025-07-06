package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitCaster_G2Kg(t *testing.T) {
	tests := []struct {
		g  float64
		kg float64
		p  int
	}{
		{1000, 1, 0},
		{1001, 1, 0},
		{1002, 1, 2},
		{1010, 1.01, 2},
		{1014, 1.01, 2},
		{1015, 1.01, 2},
		// 四舍六入五考虑，五后非零就进一，五后为零看奇偶，五前为偶应舍去，五前为奇要进一
	}
	for _, tt := range tests {
		uc := NewUnitCaster(tt.p)
		assert.Equalf(t, tt.kg, uc.G2Kg(tt.g).Float(), "G2Kg(%v)", tt.g)
	}
}
