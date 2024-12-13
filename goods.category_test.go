package temu

import (
	"github.com/hiscaler/gox/jsonx"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_goodsCategoryService_Query(t *testing.T) {
	params := GoodsCategoryQueryParams{}
	_, err := temuClient.Services.Goods.Category.Query(ctx, params)
	assert.Equalf(t, nil, err, "Services.Goods.Category.Query(ctx, %s)", jsonx.ToPrettyJson(params))
}
