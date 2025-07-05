package temu

import (
	"testing"

	"github.com/hiscaler/gox/jsonx"
	"github.com/stretchr/testify/assert"
)

func Test_goodsCertificationService_Query(t *testing.T) {
	params := GoodsCertificationQueryParams{
		ProductId: 43755443910,
	}
	_, err := temuClient.Services.Goods.Certification.Query(ctx, params)
	assert.Equalf(t, nil, err, "Services.Goods.Certification.Query(ctx, %s)", jsonx.ToPrettyJson(params))
}
