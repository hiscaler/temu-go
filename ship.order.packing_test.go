package temu

import (
	"github.com/hiscaler/gox/jsonx"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShipOrderPackingService_Match(t *testing.T) {
	params := ShipOrderPackingMatchRequest{
		DeliveryOrderSnList: []string{"FH2408231977953"},
	}
	_, err := temuClient.Services.ShipOrderPackingService.Match(params)
	assert.Nilf(t, err, "temuClient.Services.ShipOrderPackingService.Match(%s)", jsonx.ToJson(params, "{}"))
}
