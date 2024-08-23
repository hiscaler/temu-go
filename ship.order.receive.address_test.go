package temu

import (
	"github.com/hiscaler/temu-go/entity"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestShipOrderReceiveAddressService_All(t *testing.T) {
	subPurchaseOrderSnList := []string{"WB2408182975602", "WB240817842833"}
	items, err := temuClient.Services.ShipOrderReceiveAddress.All(subPurchaseOrderSnList...)
	assert.Nilf(t, err, "Services.ShipOrderReceiveAddress.All(%s)", strings.Join(subPurchaseOrderSnList, ", "))
	if len(items) != 0 {
		item := items[0]
		subPurchaseOrderSn := item.SubPurchaseOrderSnList[0]
		var d entity.ShipOrderReceiveAddress
		d, err = temuClient.Services.ShipOrderReceiveAddress.One(subPurchaseOrderSn)
		assert.Nilf(t, err, "Services.ShipOrderReceiveAddress.One(%s)", subPurchaseOrderSn)
		assert.Equalf(t, item.SubWarehouseId, d.SubWarehouseId, "Services.ShipOrderReceiveAddress.One(%s)", subPurchaseOrderSn)
		assert.Equalf(t, item.ReceiveAddressInfo, d.ReceiveAddressInfo, "Services.ShipOrderReceiveAddress.One(%s)", subPurchaseOrderSn)
		assert.Equalf(t, 1, len(d.SubPurchaseOrderSnList), "Services.ShipOrderReceiveAddress.One(%s)", subPurchaseOrderSn)
	}
}
