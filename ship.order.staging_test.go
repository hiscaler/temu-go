package temu

import (
	"fmt"
	"github.com/hiscaler/gox/jsonx"
	"github.com/hiscaler/temu-go/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShipOrderStagingService_Query(t *testing.T) {
	params := ShipOrderStagingQueryParams{SubPurchaseOrderSnList: []string{"WB2410143760309"}}
	_, _, _, _, err := temuClient.Services.ShipOrder.Staging.Query(ctx, params)
	assert.Nilf(t, err, "Services.ShipOrder.Staging.Query(ctx, %s)", jsonx.ToJson(params, "{}"))
}

func TestShipOrderStagingService_Add(t *testing.T) {
	req := ShipOrderStagingAddRequest{}
	req.JoinInfoList = []ShipOrderStagingAddInfo{
		{
			SubPurchaseOrderSn:  "1111111WB2409203163348",
			DeliveryAddressType: entity.DeliveryAddressTypeChineseMainland,
		},
	}
	reqJson := jsonx.ToJson(req, "{}")
	ok, results, err := temuClient.Services.ShipOrder.Staging.Add(ctx, req)
	fmt.Println(err.Error())
	assert.Equalf(t, len(req.JoinInfoList), len(results), "Services.ShipOrder.Staging.Add(%s)", reqJson)
	assert.Equalf(t, err == nil, ok, "Services.ShipOrder.Staging.Add(ctx, %s) ok", reqJson)
}

func TestShipOrderStagingService_One(t *testing.T) {
	qp := ShipOrderStagingQueryParams{}
	qp.SubPurchaseOrderSnList = []string{"WB2409203163348"}
	items, _, _, _, err := temuClient.Services.ShipOrder.Staging.Query(ctx, qp)
	if err == nil && len(items) != 0 {
		item := items[0]
		subPurchaseOrderSn := item.SubPurchaseOrderBasicVO.SubPurchaseOrderSn
		var d entity.ShipOrderStaging
		d, err = temuClient.Services.ShipOrder.Staging.One(ctx, subPurchaseOrderSn)
		assert.Nilf(t, err, "temuClient.Services.ShipOrder.Staging.One(ctx, %s)", subPurchaseOrderSn)
		assert.Equalf(t, item, d, "temuClient.Services.ShipOrder.Staging.One(ctx, %s)", subPurchaseOrderSn)
	}

}
