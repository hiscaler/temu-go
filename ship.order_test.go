package temu

import (
	"testing"
)

func TestShipOrderService_All(t *testing.T) {
	params := ShipOrderQueryParams{}
	params.Page = 1
	params.PageSize = 10
	params.SubPurchaseOrderSnList = []string{"WB2408152736638"}
	items, _, _, _, err := temuClient.Services.ShipOrder.All(params)
	if err != nil {
		t.Errorf("temuClient.Services.ShipOrderStaging.Companies: %s", err.Error())
	} else {
		t.Logf("items: %#v", items)
	}
}
