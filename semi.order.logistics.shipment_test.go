package temu

import (
	"fmt"
	"path"
	"strings"
	"testing"
)

func TestSemiOrderLogisticsShipmentQuery(t *testing.T) {
	productImage := "./uploads/aaa"
	productImage = path.Clean(productImage)
	if strings.HasPrefix(productImage, "/") {
		productImage = productImage[1:]
	}
	fmt.Println(productImage)
	// params := SemiOrderLogisticsShipmentQueryParams{
	// 	ParentOrderSn: "PO-211-02550200509992062",
	// 	OrderSn:       "211-02550219515432062",
	// }
	// result, err := temuClient.Services.SemiManaged.OrderLogisticsShipment.Query(ctx, params)
	// assert.Equalf(t, nil, err, "SemiManaged.OrderLogisticsShipment.Query(ctx, %#v)", params)
	// fmt.Printf("%#v\n", result)
}
