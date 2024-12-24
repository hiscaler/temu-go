package temu

import (
	"github.com/hiscaler/temu-go/entity"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4"
	"testing"
)

func Test_purchaseOrderService_Query(t *testing.T) {
	params := PurchaseOrderQueryParams{
		StatusList: []int{
			// entity.PurchaseOrderStatusWaitingMerchantReceive, // 待接单
			// entity.PurchaseOrderStatusMerchantReceived,       // 已接单/待发货
			// entity.PurchaseOrderStatusMerchantSend, // 已送货
			// entity.PurchaseOrderStatusPlatformReceived,       // 已收货
			// entity.PurchaseOrderStatusPlatformR收
			//			// entity.PurchaseOrderStatusPlatformReturned,       // 已验ejected,       // 已拒收/全部退回
			// entity.PurchaseOrderStatusPlatformApproved,       // 已验收
			// entity.PurchaseOrderStatusPlatformPutInStorage,   // 已入库
			// entity.PurchaseOrderStatusDiscard,                // 作废
			// entity.PurchaseOrderStatusTimeout,                // 已超时
			// entity.PurchaseOrderStatusCancel,                 // 已取消
		},
		SubPurchaseOrderSnList: []string{"FH2411191576410"},
		// OrderType:              null.IntFrom(entity.StockTypeCustomized),
		IsCustomGoods: null.NewBool(true, true),
		// PurchaseStockType: null.IntFrom(0),
		// JoinDeliveryPlatform: null.NewBool(false, true),
		// IsFirst:              null.NewBool(false, true),
		// UrgencyType:          null.NewInt(0, false),
	}
	params.Page = 1
	params.PageSize = 10
	items, _, err := temuClient.Services.PurchaseOrder.Query(ctx, params)
	assert.Equalf(t, nil, err, "Services.PurchaseOrder.Query(ctx, %#v) err", params)
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
