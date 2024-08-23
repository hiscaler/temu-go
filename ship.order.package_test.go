package temu

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShipOrderPackageService_One(t *testing.T) {
	deliveryOrderSn := "FH2408221920008"
	_, err := temuClient.Services.ShipOrderPackageService.One(deliveryOrderSn)
	assert.Nilf(t, err, "temuClient.Services.ShipOrderPackageService.One(%s)", deliveryOrderSn)
}
