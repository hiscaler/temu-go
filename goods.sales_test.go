package temu

import (
	"github.com/hiscaler/gox/jsonx"
	"github.com/hiscaler/temu-go/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_goodsSalesService_All(t *testing.T) {
	params := GoodsSalesQueryParams{}
	items, err := temuClient.Services.GoodsSales.All(ctx, params)
	assert.Equalf(t, nil, err, "Services.GoodsSales.All(ctx, %s)", jsonx.ToPrettyJson(params))

	if len(items) != 0 {
		item := items[0]
		var sales entity.GoodsSales
		// 根据商品 SKC ID 查询
		sales, err = temuClient.Services.GoodsSales.One(ctx, item.ProductSkcId)
		assert.Equalf(t, nil, err, "Services.GoodsSales.One(ctx, %d)", item.ProductSkcId)
		assert.Equalf(t, item, sales, "Services.GoodsSales.One(ctx, %d)", item.ProductSkcId)
	}
}
