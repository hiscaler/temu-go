package temu

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSemiOrderLogisticsShipmentQuery(t *testing.T) {
	params := SemiOrderLogisticsShipmentGetRequest{
		ParentOrderSn: "PO-211-02550200509992062",
		OrderSn:       "211-02550219515432062",
	}
	result, err := temuClient.Services.SemiManaged.OrderLogisticsShipment.Query(ctx, params)
	assert.Equalf(t, nil, err, "SemiManaged.OrderLogisticsShipment.Query(ctx, %#v)", params)
	fmt.Printf("%#v\n", result)
}
