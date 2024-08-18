package temu

import (
	"github.com/hiscaler/gox/jsonx"
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
	reqJson := jsonx.ToJson(req, "{}")
	ok, results, err := temuClient.Services.ShipOrderStaging.Add(req)
	assert.Equalf(t, len(req.JoinInfoList), len(results), "Services.ShipOrderStaging.Add(%s)", reqJson)
	assert.Equalf(t, err != nil, ok, "Services.ShipOrderStaging.Add(%s) ok", reqJson)
}

func TestShipOrderStagingService_One(t *testing.T) {
	items, _, _, _, err := temuClient.Services.ShipOrderStaging.All(ShipOrderStagingQueryParams{})
	if err == nil && len(items) != 0 {
		item := items[0]
		subPurchaseOrderSn := item.SubPurchaseOrderBasicVO.SubPurchaseOrderSn
		var d entity.ShipOrderStaging
		d, err = temuClient.Services.ShipOrderStaging.One(subPurchaseOrderSn)
		assert.Nilf(t, err, "temuClient.Services.ShipOrderStaging.One(%s)", subPurchaseOrderSn)
		assert.Equalf(t, item, d, "temuClient.Services.ShipOrderStaging.One(%s)", subPurchaseOrderSn)

	}

}
