package temu

import (
	"github.com/hiscaler/temu-go/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_purchaseOrderService_All(t *testing.T) {
	params := PurchaseOrderQueryParams{
		StatusList:                  []int{entity.PurchaseOrderStatusWaitingMerchantReceive, entity.PurchaseOrderStatusMerchantReceived},
		SubPurchaseOrderSnList:      []string{},
		OriginalPurchaseOrderSnList: []string{},
	}
	items, _, err := temuClient.Services.PurchaseOrderService.All(params)
	assert.Equalf(t, nil, err, "Services.PurchaseOrderService.All(%#v) err", params)

	if len(items) != 0 {
		item := items[0]
		var order entity.PurchaseOrder

		// 根据母订单号查询
		order, err = temuClient.Services.PurchaseOrderService.One(item.OriginalPurchaseOrderSn)
		assert.Equalf(t, nil, err, "Services.PurchaseOrderService.One(%s)", item.OriginalPurchaseOrderSn)
		assert.Equalf(t, item, order, "Services.PurchaseOrderService.One(%s)", item.OriginalPurchaseOrderSn)

		// 根据子订单号查询
		order, err = temuClient.Services.PurchaseOrderService.One(item.SubPurchaseOrderSn)
		assert.Equalf(t, nil, err, "Services.PurchaseOrderService.One(%s)", item.SubPurchaseOrderSn)
		assert.Equalf(t, item, order, "Services.PurchaseOrderService.One(%s)", item.SubPurchaseOrderSn)
	}
}

func Test_purchaseOrderService_One(t *testing.T) {
	purchaseOrderSn := "WB2408173170013"
	d, err := temuClient.Services.PurchaseOrderService.One(purchaseOrderSn)
	assert.Equalf(t, nil, err, "Services.PurchaseOrderService.One(%s) err", purchaseOrderSn)
	assert.Equalf(t, purchaseOrderSn, d.SubPurchaseOrderSn, "Services.PurchaseOrderService.One(%s) value", purchaseOrderSn)
}
