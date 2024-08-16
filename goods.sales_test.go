package temu

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_goodsSalesService_BoxMark(t *testing.T) {
	params := GoodsSalesQueryParams{}
	s, err := temuClient.Services.GoodsSalesService.All(params)
	fmt.Println(s)
	assert.Equalf(t, nil, err, "test1")
}
