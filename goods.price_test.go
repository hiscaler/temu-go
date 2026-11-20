package temu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_goodsPriceService_Query(t *testing.T) {
	items, err := temuClient.Services.Goods.Price.Query(ctx, GoodsPriceQueryParams{
		ProductSkuIds: []int64{1},
	})
	assert.NotEqual(t, err, nil)
	_ = items
}
