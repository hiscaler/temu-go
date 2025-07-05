package temu

import (
	"testing"

	"github.com/hiscaler/gox/jsonx"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4"
)

func Test_goodsCategoryService_Query(t *testing.T) {
	params := GoodsCategoryQueryParams{
		ParentCatId: null.NewInt(1, true),
	}
	_, err := temuClient.Services.Goods.Category.Query(ctx, params)
	assert.Equalf(t, nil, err, "Services.Goods.Category.Query(ctx, %s)", jsonx.ToPrettyJson(params))
}
