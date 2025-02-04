package temu

import (
	"github.com/hiscaler/temu-go/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_mallDeliveryAddressService_Query(t *testing.T) {
	addresses, err := temuClient.Services.Mall.DeliveryAddress.Query(ctx)
	assert.Equal(t, nil, err, "Services.Mall.DeliveryAddress.Query(ctx)")

	for _, address := range addresses {
		var addr entity.DeliveryAddress
		addr, err = temuClient.Services.Mall.DeliveryAddress.One(ctx, address.ID)
		assert.Equalf(t, nil, err, "Services.Mall.DeliveryAddress.One(ctx, %d)", address.ID)
		assert.Equalf(t, addr, address, "Services.Mall.DeliveryAddress.One(ctx, %d)", address.ID)
	}
}
