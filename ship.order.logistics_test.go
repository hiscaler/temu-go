package temu

import (
	"fmt"
	"github.com/hiscaler/temu-go/entity"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4"
	"testing"
)

func Test_shipOrderLogisticsService_Match(t *testing.T) {
	purchaseOrder, err := temuClient.Services.PurchaseOrder.One(ctx, "WB2412251860849")
	assert.Equal(t, nil, err, "PurchaseOrder.One")
	var receiveAddress entity.ShipOrderReceiveAddress
	receiveAddress, err = temuClient.Services.ShipOrder.ReceiveAddress.One(ctx, purchaseOrder.SubPurchaseOrderSn)
	assert.Equal(t, nil, err, "ShipOrder.ReceiveAddress.One")
	shipOrders, _, _, _, err := temuClient.Services.ShipOrder.Query(ctx, ShipOrderQueryParams{
		SubPurchaseOrderSnList: []string{purchaseOrder.SubPurchaseOrderSn},
	})
	assert.Equal(t, nil, err, "ShipOrder.Query")
	// assert.Equal(t, 1, err, "ShipOrder.Query result")

	shipOrder := shipOrders[0]
	req := LogisticsMatchRequest{
		DeliveryAddressId:         12,
		PredictTotalPackageWeight: 1000,
		UrgencyType:               null.IntFrom(1),
		SubWarehouseId:            shipOrder.SubWarehouseId,
		TotalPackageNum:           2,
		ReceiveAddressInfo:        &receiveAddress.ReceiveAddressInfo,
		DeliveryOrderSns:          []string{shipOrder.DeliveryOrderSn},
	}
	items, err := temuClient.Services.ShipOrder.Logistics.Match(ctx, req)
	assert.Equal(t, nil, err, "Services.Logistics.Match(ctx)")
	for _, item := range items {
		fmt.Printf("%#v", item)
	}
}
