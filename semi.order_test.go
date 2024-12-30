package temu

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOrder(t *testing.T) {
	params := OrderQueryParams{
		ParentOrderStatus: 2,
	}
	params.Page = 1
	params.PageSize = 10
	items, total, totalPage, isLatestPage, err := temuClient.Services.SemiManaged.Order.Query(ctx, params)
	assert.Equalf(t, nil, err, "Services.SemiManaged.Order.Query(ctx, %#v) err", params)
	println(items, total, totalPage, isLatestPage)
}
