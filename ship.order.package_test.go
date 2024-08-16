package temu

import (
	"testing"
)

func TestShipOrderPackageService_One(t *testing.T) {
	data, err := temuClient.Services.ShipOrderPackageService.One("FH2408131785686")
	if err != nil {
		t.Errorf("temuClient.Services.ShipOrderStaging.Companies: %s", err.Error())
	} else {
		t.Logf("items: %#v", data)
	}
}
