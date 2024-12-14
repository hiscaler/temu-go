package temu

import (
	"fmt"
	"github.com/hiscaler/temu-go/entity"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4"
	"testing"
)

func Test_logisticsCompanyService_Companies(t *testing.T) {
	companies, err := temuClient.Services.Logistics.Companies(ctx)
	assert.Equal(t, nil, err, "Services.Logistics.Companies(ctx)")
	for _, company := range companies {
		var com entity.LogisticsShipmentCompany
		com, err = temuClient.Services.Logistics.Company(ctx, company.ShipId)
		assert.Equalf(t, nil, err, "Services.Logistics.Company(ctx, %d)", company.ShipId)
		assert.Equalf(t, company, com, "Services.Logistics.Company(ctx, %d)", company.ShipId)
	}
}

func Test_logisticsService_Match(t *testing.T) {
	purchaseOrder, err := temuClient.Services.PurchaseOrder.One(ctx, "WB2410011321611")
	assert.Equal(t, nil, err, "PurchaseOrder.One")
	var receiveAddress entity.ShipOrderReceiveAddress
	receiveAddress, err = temuClient.Services.ShipOrderReceiveAddress.One(ctx, purchaseOrder.SubPurchaseOrderSn)
	assert.Equal(t, nil, err, "ShipOrderReceiveAddress.One")
	shipOrders, _, _, _, err := temuClient.Services.ShipOrder.Query(ctx, ShipOrderQueryParams{
		SubPurchaseOrderSnList: []string{purchaseOrder.SubPurchaseOrderSn},
	})
	assert.Equal(t, nil, err, "ShipOrder.Query")
	// assert.Equal(t, 1, err, "ShipOrder.Query result")

	shipOrder := shipOrders[0]
	req := LogisticsMatchRequest{
		DeliveryAddressId:         shipOrder.DeliveryAddressId.Int64,
		PredictTotalPackageWeight: 1000,
		UrgencyType:               null.IntFrom(1),
		SubWarehouseId:            shipOrder.SubWarehouseId,
		QueryStandbyExpress:       null.BoolFrom(false),
		TotalPackageNum:           2,
		ReceiveAddressInfo:        receiveAddress.ReceiveAddressInfo,
		DeliveryOrderSns:          []string{shipOrder.DeliveryOrderSn},
	}
	items, err := temuClient.Services.Logistics.Match(ctx, req)
	assert.Equal(t, nil, err, "Services.Logistics.Match(ctx)")
	for _, item := range items {
		fmt.Printf("%#v", item)
	}
}
