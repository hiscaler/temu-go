package temu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrder(t *testing.T) {
	params := OrderQueryParams{
		// ParentOrderStatus: 2,
		RegionId:          211,
		ParentOrderSnList: []string{"PO-211-05430492122150798"},
	}
	params.Page = 1
	params.PageSize = 10
	items, total, totalPage, isLatestPage, err := temuClient.Services.SemiManaged.Order.Query(ctx, params)
	assert.Equalf(t, nil, err, "Services.SemiManaged.Order.Query(ctx, %#v) err", params)
	println(items, total, totalPage, isLatestPage)
}
