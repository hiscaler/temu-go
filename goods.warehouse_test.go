package temu

import (
	"github.com/hiscaler/gox/jsonx"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_goodsWarehouseService_Query(t *testing.T) {
	params := GoodsWarehouseQueryParams{SiteIdList: []int64{6545025}}
	err := temuClient.Services.Goods.Warehouse.Query(ctx, params)
	assert.Equalf(t, nil, err, "Services.Goods.Warehouse.Query(ctx, %s)", jsonx.ToPrettyJson(params))
}
