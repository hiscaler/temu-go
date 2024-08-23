package temu

import (
	"github.com/hiscaler/temu-go/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_mallAddressService_All(t *testing.T) {
	addresses, err := temuClient.Services.MallAddress.All()
	assert.Equal(t, nil, err, "Test_mallAddressService_All")

	for _, address := range addresses {
		var addr entity.MallAddress
		addr, err = temuClient.Services.MallAddress.One(address.ID)
		assert.Equalf(t, nil, err, "Services.MallAddress.One(%d)", address.ID)
		assert.Equalf(t, addr, address, "Services.MallAddress.One(%d)", address.ID)
	}
}
