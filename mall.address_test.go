package temu

import (
	"github.com/hiscaler/temu-go/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_mallAddressService_Query(t *testing.T) {
	addresses, err := temuClient.Services.Mall.Address.Query(ctx)
	assert.Equal(t, nil, err, "Test_mallAddressService_Query")

	for _, address := range addresses {
		var addr entity.MallAddress
		addr, err = temuClient.Services.Mall.Address.One(ctx, address.ID)
		assert.Equalf(t, nil, err, "Services.Mall.Address.One(ctx, %d)", address.ID)
		assert.Equalf(t, addr, address, "Services.Mall.Address.One(ctx, %d)", address.ID)
	}
}
