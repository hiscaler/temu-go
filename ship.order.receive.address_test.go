package temu

import (
	"fmt"
	"testing"
)

func TestShipOrderReceiveAddressService_All(t *testing.T) {
	items, err := temuClient.Services.ShipOrderReceiveAddressService.All("3243242")
	if err != nil {
		t.Errorf("temuClient.Services.ShipOrderReceiveAddressService.All: %s", err.Error())
	} else {
		fmt.Printf("items: %#v", items)
	}
}
