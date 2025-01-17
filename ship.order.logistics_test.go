package temu

import (
	"fmt"
	"github.com/hiscaler/temu-go/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_shipOrderLogisticsService_Match(t *testing.T) {
	purchaseOrder, err := temuClient.Services.PurchaseOrder.One(ctx, "WB2501082789597")
	assert.Equal(t, nil, err, "PurchaseOrder.One")
	var receiveAddress entity.ShipOrderReceiveAddress
	receiveAddress, err = temuClient.Services.ShipOrder.ReceiveAddress.One(ctx, purchaseOrder.SubPurchaseOrderSn)
	assert.Equal(t, nil, err, "ShipOrder.ReceiveAddress.One")
	shipOrders, _, _, _, err := temuClient.Services.ShipOrder.Query(ctx, ShipOrderQueryParams{
		SubPurchaseOrderSnList: []string{purchaseOrder.SubPurchaseOrderSn},
	})
	assert.Equal(t, nil, err, "ShipOrder.Query")
	// assert.Equal(t, 1, err, "ShipOrder.Query result")

	if len(shipOrders) != 0 {
		shipOrder := shipOrders[0]
		req := LogisticsMatchRequest{
			DeliveryAddressId:         20172531361454,
			PredictTotalPackageWeight: 1000,
			// UrgencyType:               null.IntFrom(1),
			SubWarehouseId:     shipOrder.SubWarehouseId,
			TotalPackageNum:    len(shipOrder.PackageList),
			ReceiveAddressInfo: &receiveAddress.ReceiveAddressInfo,
			DeliveryOrderSns:   []string{shipOrder.DeliveryOrderSn},
		}
		items, err := temuClient.Services.ShipOrder.Logistics.Match(ctx, req)
		assert.Equal(t, nil, err, "Services.Logistics.Match(ctx)")
		for _, item := range items {
			fmt.Printf("%#v\n", item)
		}
	}
}
