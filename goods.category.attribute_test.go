package temu

import (
	"fmt"
	"testing"

	"github.com/hiscaler/gox/jsonx"
	"github.com/stretchr/testify/assert"
)

func Test_goodsCategoryAttributeService_Query(t *testing.T) {
	var categoryId int64 = 10023
	_, err := temuClient.Services.Goods.Category.Attribute.Query(ctx, categoryId)
	assert.Equalf(t, nil, err, "Services.Goods.Category.Attribute.Query(ctx, %s)", jsonx.ToPrettyJson(categoryId))
	fmt.Println(err)
}
