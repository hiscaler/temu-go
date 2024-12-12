package temu

import (
	"testing"
)

func Test_goodsWarehouseService_Query(t *testing.T) {
	// params := GoodsSalesQueryParams{}
	_ = temuClient.Services.Goods.Warehouse.Query(ctx)
	// assert.Equalf(t, nil, err, "Services.GoodsSales.Query(ctx, %s)", jsonx.ToPrettyJson(params))
	//
	// if len(items) != 0 {
	// 	item := items[0]
	// 	var sales entity.GoodsSales
	// 	// 根据商品 SKC ID 查询
	// 	sales, err = temuClient.Services.Goods.Sales.One(ctx, item.ProductSkcId)
	// 	assert.Equalf(t, nil, err, "Services.GoodsSales.One(ctx, %d)", item.ProductSkcId)
	// 	assert.Equalf(t, item, sales, "Services.GoodsSales.One(ctx, %d)", item.ProductSkcId)
	// }
}
