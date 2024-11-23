package temu

import (
	"github.com/hiscaler/temu-go/entity"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestShipOrderReceiveAddressService_Query(t *testing.T) {
	subPurchaseOrderSnList := []string{"WB2411181848215"}
	items, err := temuClient.Services.ShipOrderReceiveAddress.Query(ctx, subPurchaseOrderSnList...)
	assert.Nilf(t, err, "Services.ShipOrderReceiveAddress.Query(ctx, %s)", strings.Join(subPurchaseOrderSnList, ", "))
	if len(items) != 0 {
		item := items[0]
		subPurchaseOrderSn := item.SubPurchaseOrderSnList[0]
		var d entity.ShipOrderReceiveAddress
		d, err = temuClient.Services.ShipOrderReceiveAddress.One(ctx, subPurchaseOrderSn)
		assert.Nilf(t, err, "Services.ShipOrderReceiveAddress.One(ctx, %s)", subPurchaseOrderSn)
		assert.Equalf(t, item.SubWarehouseId, d.SubWarehouseId, "Services.ShipOrderReceiveAddress.One(%s)", subPurchaseOrderSn)
		assert.Equalf(t, item.ReceiveAddressInfo, d.ReceiveAddressInfo, "Services.ShipOrderReceiveAddress.One(%s)", subPurchaseOrderSn)
		assert.Equalf(t, 1, len(d.SubPurchaseOrderSnList), "Services.ShipOrderReceiveAddress.One(%s)", subPurchaseOrderSn)
	}
}
