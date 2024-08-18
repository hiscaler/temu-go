package temu

import (
	"github.com/hiscaler/temu-go/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShipOrderStagingService_All(t *testing.T) {
	params := ShipOrderStagingQueryParams{}
	_, _, _, _, err := temuClient.Services.ShipOrderStaging.All(params)
	assert.Nilf(t, err, "Services.ShipOrderStaging.All(%#v)", params)
}

func TestShipOrderStagingService_Add(t *testing.T) {
	req := ShipOrderStagingAddRequest{}
	req.JoinInfoList = []ShipOrderStagingAddInfo{
		{
			SubPurchaseOrderSn:  "WB2408173170013",
			DeliveryAddressType: entity.DeliveryAddressTypeChineseMainland,
		},
	}
	_, results, _ := temuClient.Services.ShipOrderStaging.Add(req)
	assert.Equalf(t, len(req.JoinInfoList), len(results), "Services.ShipOrderStaging.Add(%#v)", req)
}

func TestShipOrderStagingService_One(t *testing.T) {
	_, err := temuClient.Services.ShipOrderStaging.One("WB2408163258440")
	if err != nil {
		t.Errorf("temuClient.Services.ShipOrderStaging.One: %s", err.Error())
	}
}
