package temu

import (
	"github.com/hiscaler/temu-go/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_purchaseOrderService_All(t *testing.T) {
	params := PurchaseOrderQueryParams{
		StatusList: []int{
			// entity.PurchaseOrderStatusWaitingMerchantReceive, // 待接单
			// entity.PurchaseOrderStatusMerchantReceived,       // 已接单/待发货
			// entity.PurchaseOrderStatusMerchantSend, // 已送货
			// entity.PurchaseOrderStatusPlatformReceived,       // 已收货
			// entity.PurchaseOrderStatusPlatformRejected,       // 已拒收
			// entity.PurchaseOrderStatusPlatformReturned,       // 已验收/全部退回
			// entity.PurchaseOrderStatusPlatformApproved,       // 已验收
			// entity.PurchaseOrderStatusPlatformPutInStorage,   // 已入库
			// entity.PurchaseOrderStatusDiscard,                // 作废
			// entity.PurchaseOrderStatusTimeout,                // 已超时
			// entity.PurchaseOrderStatusCancel,                 // 已取消
		},
		SubPurchaseOrderSnList: []string{"WB2408222923964"},
		// OriginalPurchaseOrderSnList: []string{},
		// IsCustomGoods:               true,
		// JoinDeliveryPlatform:        true,
	}
	params.PageSize = 1000
	items, _, err := temuClient.Services.PurchaseOrder.All(ctx, params)
	assert.Equalf(t, nil, err, "Services.PurchaseOrder.All(ctx, %#v) err", params)

	if len(items) != 0 {
		item := items[0]
		var order entity.PurchaseOrder

		// 根据母订单号查询
		order, err = temuClient.Services.PurchaseOrder.One(ctx, item.OriginalPurchaseOrderSn)
		assert.Equalf(t, nil, err, "Services.PurchaseOrder.One(ctx, %s)", item.OriginalPurchaseOrderSn)
		assert.Equalf(t, item, order, "Services.PurchaseOrder.One(ctx, %s)", item.OriginalPurchaseOrderSn)

		// 根据子订单号查询
		order, err = temuClient.Services.PurchaseOrder.One(ctx, item.SubPurchaseOrderSn)
		assert.Equalf(t, nil, err, "Services.PurchaseOrder.One(ctx, %s)", item.SubPurchaseOrderSn)
		assert.Equalf(t, item, order, "Services.PurchaseOrder.One(ctx, %s)", item.SubPurchaseOrderSn)
	}
}

func Test_purchaseOrderService_One(t *testing.T) {
	purchaseOrderSn := "WB2408173170013"
	d, err := temuClient.Services.PurchaseOrder.One(ctx, purchaseOrderSn)
	assert.Equalf(t, nil, err, "Services.PurchaseOrder.One(ctx, %s) err", purchaseOrderSn)
	assert.Equalf(t, purchaseOrderSn, d.SubPurchaseOrderSn, "Services.PurchaseOrder.One(ctx, %s) value", purchaseOrderSn)
}
