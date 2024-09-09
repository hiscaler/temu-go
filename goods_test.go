package temu

import (
	"github.com/hiscaler/gox/jsonx"
	"github.com/hiscaler/temu-go/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_goodsService_All(t *testing.T) {
	params := GoodsQueryParams{}
	params.PageSize = 2
	items, _, _, _, err := temuClient.Services.Goods.All(ctx, params)
	assert.Equalf(t, nil, err, "Services.Goods.All(ctx, %s)", jsonx.ToPrettyJson(params))
	_ = items
	if len(items) != 0 {
		item := items[0]
		var sales entity.Goods
		// 根据商品 SKC ID 查询
		sales, err = temuClient.Services.Goods.One(ctx, item.ProductSkcId)
		assert.Equalf(t, nil, err, "Services.Goods.One(ctx, %d)", item.ProductSkcId)
		assert.Equalf(t, item, sales, "Services.Goods.One(ctx, %d)", item.ProductSkcId)
	}
}
