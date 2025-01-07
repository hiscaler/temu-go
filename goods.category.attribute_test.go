package temu

import (
	"github.com/hiscaler/gox/jsonx"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_goodsCategoryAttributeService_Query(t *testing.T) {
	var categoryId int64 = 653
	_, err := temuClient.Services.Goods.Category.Attribute.Query(ctx, categoryId)
	assert.Equalf(t, nil, err, "Services.Goods.Category.Attribute.Query(ctx, %s)", jsonx.ToPrettyJson(categoryId))
}
