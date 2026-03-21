package temu

import (
	"testing"

	"github.com/hiscaler/temu-go/entity"
	"github.com/stretchr/testify/assert"
)

func TestOrder(t *testing.T) {
	params := SemiOrderQueryParams{
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
// TestIs 测试标签是否包含指定标签
func TestIs(t *testing.T) {
	labels := entity.Labels{
		entity.Label{
			Name:  "A",
			Value: 1,
		},
		entity.Label{
			Name:  "B",
			Value: 0,
		},
		entity.Label{
			Name:  "C",
			Value: 1,
		},
	}
	assert.Equalf(t, true, labels.Is("A", "C"), "Is(%v, %v)", "A", "C")
	assert.Equalf(t, false, labels.Is("A", "B"), "Is(%v, %v)", "A", "B")
}
