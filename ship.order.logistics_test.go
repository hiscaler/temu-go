package temu

import (
	"fmt"
	"testing"

	"github.com/hiscaler/temu-go/entity"
	"github.com/stretchr/testify/assert"
)

func Test_shipOrderLogisticsService_Match(t *testing.T) {
	number := "WB2501082789597"
	purchaseOrder, err := temuClient.Services.PurchaseOrder.One(ctx, number)
	assert.Equalf(t, nil, err, "PurchaseOrder.One(ctx, %s)", number)

	// 发货地址
	var deliveryAddress entity.DeliveryAddress
	deliveryAddresses, err := temuClient.Services.Mall.DeliveryAddress.Query(ctx)
	assert.Equal(t, nil, err, "Mall.DeliveryAddress.One(ctx)")
	if len(deliveryAddresses) != 0 {
		deliveryAddress = deliveryAddresses[0]
	}

	// 收货地址
	var receiveAddress entity.ShipOrderReceiveAddress
	receiveAddress, err = temuClient.Services.ShipOrder.ReceiveAddress.One(ctx, purchaseOrder.SubPurchaseOrderSn)
	assert.Equal(t, nil, err, "ShipOrder.ReceiveAddress.One")
	params := ShipOrderQueryParams{
		SubPurchaseOrderSnList: []string{purchaseOrder.SubPurchaseOrderSn},
	}
	shipOrders, _, _, _, err := temuClient.Services.ShipOrder.Query(ctx, params)
	assert.Equalf(t, nil, err, "ShipOrder.Query(ctx, %#v)", params)

	if len(shipOrders) != 0 {
		shipOrder := shipOrders[0]
		req := LogisticsMatchRequest{
			DeliveryAddressId:         deliveryAddress.ID,
			PredictTotalPackageWeight: 1000,
			// UrgencyType:               null.IntFrom(1),
			SubWarehouseId:     shipOrder.SubWarehouseId,
			TotalPackageNum:    len(shipOrder.PackageList),
			ReceiveAddressInfo: &receiveAddress.ReceiveAddressInfo,
			DeliveryOrderSns:   []string{shipOrder.DeliveryOrderSn},
		}
		items, err := temuClient.Services.ShipOrder.Logistics.Match(ctx, req)
		assert.Equalf(t, nil, err, "Services.Logistics.Match(ctx, %#v)", req)
		for _, item := range items {
			fmt.Printf("%#v\n", item)
		}
	}
}
