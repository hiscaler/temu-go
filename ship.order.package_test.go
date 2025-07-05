package temu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShipOrderPackageService_One(t *testing.T) {
	deliveryOrderSn := "FH2410242429251"
	_, err := temuClient.Services.ShipOrder.Package.One(ctx, deliveryOrderSn)
	assert.Nilf(t, err, "temuClient.Services.ShipOrder.Package.One(ctx, %s)", deliveryOrderSn)
}
