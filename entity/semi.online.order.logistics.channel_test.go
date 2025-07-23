package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSemiOnlineOrderLogisticsChannel_ParseEstimatedAmount(t *testing.T) {
	tests := map[string]float64{
		"1-1":     0,
		"$91.21":  91.21,
		"$ 91.21": 91.21,
		"$ 91":    91,
		"$ .01":   0.01,
		"$ 1.234": 1.234,
	}
	for text, amount := range tests {
		d := SemiOnlineOrderLogisticsChannel{EstimatedAmount: text}
		v, _ := d.ParseEstimatedAmount()
		assert.Equalf(t, amount, v, "Amount(%s)", text)
	}
}

func TestSemiOnlineOrderLogisticsChannel_DeliveryDays(t *testing.T) {
	tests := map[string][]float64{
		"1-1":                     {1, 1},
		"1  -2":                   {1, 2},
		"预估$91.21; USD; 1-2工作日送达": {1, 2},
		"预估$91.21; USD; 1 - 2工作日送达":     {1, 2},
		"预估$91.21; USD; 1.0 - 2.5工作日送达": {1, 2.5},
	}
	for text, days := range tests {
		d := SemiOnlineOrderLogisticsChannel{EstimatedText: text}
		minDays, maxDays, _ := d.DeliveryDays()
		assert.Equalf(t, days[0], minDays, "DeliveryDays(%s)", text)
		assert.Equalf(t, days[1], maxDays, "DeliveryDays(%s)", text)
	}
}
