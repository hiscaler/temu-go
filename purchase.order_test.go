package temu

import (
	"fmt"
	"testing"
)

func Test_purchaseOrderService_All(t *testing.T) {
	params := PurchaseOrderQueryParams{
		Page:                   1,
		PageSize:               10,
		SettlementType:         1,
		SubPurchaseOrderSnList: []string{"WB2408152983266"},
	}
	s, err := temuClient.Services.PurchaseOrderService.All(params)
	fmt.Println(s)
	fmt.Println(err)
	// assert.Equalf(t, nil, err, "test1")
}
