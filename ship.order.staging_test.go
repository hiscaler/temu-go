package temu

import (
	"github.com/hiscaler/gox/jsonx"
	"github.com/hiscaler/temu-go/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShipOrderStagingService_All(t *testing.T) {
	params := ShipOrderStagingQueryParams{SubPurchaseOrderSnList: []string{"WB2408163594797"}}
	_, _, _, _, err := temuClient.Services.ShipOrderStaging.All(ctx, params)
	assert.Nilf(t, err, "Services.ShipOrderStaging.All(ctx, %s)", jsonx.ToJson(params, "{}"))
}

func TestShipOrderStagingService_Add(t *testing.T) {
	req := ShipOrderStagingAddRequest{}
	req.JoinInfoList = []ShipOrderStagingAddInfo{
		{
			SubPurchaseOrderSn:  "WB2408182975602",
			DeliveryAddressType: entity.DeliveryAddressTypeChineseMainland,
		},
	}
	reqJson := jsonx.ToJson(req, "{}")
	ok, results, err := temuClient.Services.ShipOrderStaging.Add(ctx, req)
	assert.Equalf(t, len(req.JoinInfoList), len(results), "Services.ShipOrderStaging.Add(%s)", reqJson)
	assert.Equalf(t, err == nil, ok, "Services.ShipOrderStaging.Add(ctx, %s) ok", reqJson)
}

func TestShipOrderStagingService_One(t *testing.T) {
	items, _, _, _, err := temuClient.Services.ShipOrderStaging.All(ctx, ShipOrderStagingQueryParams{})
	if err == nil && len(items) != 0 {
		item := items[0]
		subPurchaseOrderSn := item.SubPurchaseOrderBasicVO.SubPurchaseOrderSn
		var d entity.ShipOrderStaging
		d, err = temuClient.Services.ShipOrderStaging.One(ctx, subPurchaseOrderSn)
		assert.Nilf(t, err, "temuClient.Services.ShipOrderStaging.One(ctx, %s)", subPurchaseOrderSn)
		assert.Equalf(t, item, d, "temuClient.Services.ShipOrderStaging.One(ctx, %s)", subPurchaseOrderSn)
	}

}
