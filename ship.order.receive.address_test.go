package temu

import (
	"github.com/hiscaler/temu-go/entity"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestShipOrderReceiveAddressService_Query(t *testing.T) {
	subPurchaseOrderSnList := []string{"WB2411181848215"}
	items, err := temuClient.Services.ShipOrder.ReceiveAddress.Query(ctx, subPurchaseOrderSnList...)
	assert.Nilf(t, err, "Services.ShipOrder.ReceiveAddress.Query(ctx, %s)", strings.Join(subPurchaseOrderSnList, ", "))
	if len(items) != 0 {
		item := items[0]
		subPurchaseOrderSn := item.SubPurchaseOrderSnList[0]
		var d entity.ShipOrderReceiveAddress
		d, err = temuClient.Services.ShipOrder.ReceiveAddress.One(ctx, subPurchaseOrderSn)
		assert.Nilf(t, err, "Services.ShipOrder.ReceiveAddress.One(ctx, %s)", subPurchaseOrderSn)
		assert.Equalf(t, item.SubWarehouseId, d.SubWarehouseId, "Services.ShipOrder.ReceiveAddress.One(%s)", subPurchaseOrderSn)
		assert.Equalf(t, item.ReceiveAddressInfo, d.ReceiveAddressInfo, "Services.ShipOrder.ReceiveAddress.One(%s)", subPurchaseOrderSn)
		assert.Equalf(t, 1, len(d.SubPurchaseOrderSnList), "Services.ShipOrder.ReceiveAddress.One(%s)", subPurchaseOrderSn)
	}
}
