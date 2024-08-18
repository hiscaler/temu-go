package temu

import (
	"github.com/hiscaler/temu-go/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_mallAddressService_All(t *testing.T) {
	addresses, err := temuClient.Services.MallAddressService.All()
	assert.Equal(t, nil, err, "Test_mallAddressService_All")

	for _, address := range addresses {
		var addr entity.MallAddress
		addr, err = temuClient.Services.MallAddressService.One(address.ID)
		assert.Equalf(t, nil, err, "Services.MallAddressService.One(%d)", address.ID)
		assert.Equalf(t, addr, address, "Services.MallAddressService.One(%d)", address.ID)
	}
}
