package temu

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_goodsLifeCycleService_Query(t *testing.T) {
	gotItems, gotTotal, gotTotalPages, gotIsLastPage, err := temuClient.Services.Goods.LifeCycle.Query(context.Background(), GoodsLifeCycleQueryParams{
		ProductSkuIdList: []int64{3544619019},
	})

	assert.NoError(t, err)
	if err != nil {
		return
	}

	_ = gotItems
	_ = gotTotal
	assert.Equalf(t, 1, gotTotalPages, "Query()")
	assert.Equalf(t, true, gotIsLastPage, "Query()")
}
