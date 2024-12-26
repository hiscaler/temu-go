package temu

import (
	"testing"

	"github.com/hiscaler/gox/jsonx"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4"
)

func Test_goodsQuantityGet(t *testing.T) {
	params := GoodsQuantityQueryParams{
		ProductSkcId: 2646847407,
	}
	_, err := temuClient.Services.Goods.Quantity.Query(ctx, params)
	assert.Equalf(t, nil, err, "Services.Goods.Quantity.Query(ctx, %s)", jsonx.ToPrettyJson(params))
}

func Test_goodsQuantityUpdate(t *testing.T) {
	params := GoodsQuantityUpdateParams{
		QuantityChangeMode: 2,
		SkuStockChangeList: []StockChangeItem{
			{
				ProductSkuId:         9437668608,
				TargetStockAvailable: null.IntFrom(0),
				WarehouseId:          null.StringFrom("WH-02048540446873850"),
			},
		},
	}

	_, err := temuClient.Services.Goods.Quantity.Update(ctx, params)
	assert.Equalf(t, nil, err, "Services.Goods.Quantity.Update(ctx, %s)", jsonx.ToPrettyJson(params))
}
