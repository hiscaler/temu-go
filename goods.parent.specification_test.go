package temu

import (
	"fmt"
	"github.com/hiscaler/gox/jsonx"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_goodsParentSpecificationService_Query(t *testing.T) {
	parentSpecifications, err := temuClient.Services.Goods.ParentSpecification.Query(ctx)
	assert.Equalf(t, nil, err, "Services.Goods.ParentSpecification.Query(ctx)")
	fmt.Println(jsonx.ToPrettyJson(parentSpecifications))
}
