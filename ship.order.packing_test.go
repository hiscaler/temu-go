package temu

import (
	"github.com/hiscaler/gox/jsonx"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShipOrderPackingService_Match(t *testing.T) {
	req := ShipOrderPackingMatchRequest{
		DeliveryOrderSnList: []string{"FH2408231977953"},
	}
	_, err := temuClient.Services.ShipOrderPackingService.Match(req)
	assert.Nilf(t, err, "temuClient.Services.ShipOrderPackingService.Match(%s)", jsonx.ToJson(req, "{}"))
}

func TestShipOrderPackingService_Send(t *testing.T) {
	addresses, err := temuClient.Services.MallAddressService.All()
	assert.Nilf(t, err, "temuClient.Services.MallAddressService.All(): error")
	assert.Equal(t, true, len(addresses) > 0, "temuClient.Services.MallAddressService.All(): results")
	address := addresses[0]

	req := ShipOrderPackingSendRequest{
		DeliveryAddressId:   address.ID,
		DeliverMethod:       3,
		DeliveryOrderSnList: []string{"FH2408231977953"},
	}
	_, err = temuClient.Services.ShipOrderPackingService.Send(req)
	assert.Nilf(t, err, "temuClient.Services.ShipOrderPackingService.Match(%s)", jsonx.ToJson(req, "{}"))
}
