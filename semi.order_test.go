package temu

import (
	"fmt"
	"testing"

	"github.com/hiscaler/temu-go/entity"
	"github.com/stretchr/testify/assert"
)

func TestOrder(t *testing.T) {
	params := SemiOrderQueryParams{
		// ParentOrderStatus: 2,
		RegionId:          211,
		ParentOrderSnList: []string{"PO-211-13028093978231739"},
	}
	params.Page = 1
	params.PageSize = 10
	items, total, totalPage, isLatestPage, err := temuClient.Services.SemiManaged.Order.Query(ctx, params)
	assert.Equalf(t, true, items[0].OrderList[0].OrderLabel.Is("Y2_advance_sale"), "Is Y2")
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

func Test_semiOrderService_ShippingInformation(t *testing.T) {
	tests := []struct {
		name              string
		parentOrderNumber string
		wantErr           assert.ErrorAssertionFunc
	}{
		{
			"t1",
			"a",
			func(t assert.TestingT, err error, i ...interface{}) bool {
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := temuClient.Services.SemiManaged.Order.ShippingInformation(ctx, tt.parentOrderNumber)
			if !tt.wantErr(t, err, fmt.Sprintf("ShippingInformation(%v, %v)", ctx, tt.parentOrderNumber)) {
				return
			}
			assert.Equalf(t, len(got.ReceiptName) != 0, true, "ShippingInformation(%v, %v)", ctx, tt.parentOrderNumber)
		})
	}
}
