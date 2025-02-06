package temu

import (
	"github.com/hiscaler/gox/jsonx"
	"github.com/hiscaler/temu-go/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_goodsService_Query(t *testing.T) {
	params := GoodsQueryParams{
		// ProductSkcIds: []int64{2646847407},
		SkuExtCodes:    []string{"8502937482"},
		ProductSkcIds:  []int64{7469668867},
		CreatedAtStart: "2024-11-18 12:00:00",
		CreatedAtEnd:   "2024-11-18 23:59:59",
	}
	params.Page = 1
	params.PageSize = 2
	items, _, _, _, err := temuClient.Services.Goods.Query(ctx, params)
	assert.Equalf(t, nil, err, "Services.Goods.Query(ctx, %s)", jsonx.ToPrettyJson(params))
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

func Test_goodsService_Detail(t *testing.T) {
	var productId int64 = 141911679
	detail, err := temuClient.Services.Goods.Detail(ctx, productId)
	assert.Equalf(t, nil, err, "Services.Goods.One(ctx, %d)", productId)
	assert.Equalf(t, detail.ProductId, productId, "Services.Goods.One(ctx, %d)", productId)
}
