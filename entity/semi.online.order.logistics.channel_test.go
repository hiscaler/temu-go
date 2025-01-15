package entity

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSemiOnlineOrderLogisticsChannel_Amount(t *testing.T) {
	tests := map[string]float64{
		"1-1":     0,
		"$91.21":  91.21,
		"$ 91.21": 91.21,
		"$ 91":    91,
		"$ .01":   0.01,
		"$ 1.234": 1.234,
	}
	for amount, value := range tests {
		d := SemiOnlineOrderLogisticsChannel{EstimatedAmount: amount}
		v, _ := d.Amount()
		assert.Equalf(t, value, v, "Amount(%s)", amount)
	}
}

func TestSemiOnlineOrderLogisticsChannel_Timeline(t *testing.T) {
	tests := map[string][]int{
		"1-1":                              {1, 1},
		"1  -2":                            {1, 2},
		"预估$91.21; USD; 1-2工作日送达":   {1, 2},
		"预估$91.21; USD; 1 - 2工作日送达": {1, 2},
	}
	for amount, value := range tests {
		d := SemiOnlineOrderLogisticsChannel{EstimatedText: amount}
		minDays, maxDays, _ := d.Timeline()
		assert.Equalf(t, value[0], minDays, "Amount(%s)", amount)
		assert.Equalf(t, value[1], maxDays, "Amount(%s)", amount)
	}
}
