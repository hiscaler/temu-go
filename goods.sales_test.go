package temu

import (
	"testing"

	"github.com/hiscaler/gox/jsonx"
	"github.com/hiscaler/temu-go/entity"
	"github.com/stretchr/testify/assert"
)

func Test_goodsSalesService_Query(t *testing.T) {
	params := GoodsSalesQueryParams{}
	items, err := temuClient.Services.Goods.Sales.Query(ctx, params)
	assert.Equalf(t, nil, err, "Services.Goods.Sales.Query(ctx, %s)", jsonx.ToPrettyJson(params))

	if len(items) != 0 {
		item := items[0]
		var sales entity.GoodsSales
		// 根据商品 SKC ID 查询
		sales, err = temuClient.Services.Goods.Sales.One(ctx, item.ProductSkcId)
		assert.Equalf(t, nil, err, "Services.Goods.Sales.One(ctx, %d)", item.ProductSkcId)
		assert.Equalf(t, item, sales, "Services.Goods.Sales.One(ctx, %d)", item.ProductSkcId)
	}
}
