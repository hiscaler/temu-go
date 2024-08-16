package temu

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_barcodeService_BoxMark(t *testing.T) {
	params := BoxMarkBarcodeQueryParams{
		ReturnDataKey:       true,
		DeliveryOrderSnList: []string{"FH2408151907154"},
	}
	s, err := temuClient.Services.BarcodeService.BoxMark(params)
	fmt.Println(s)
	assert.Equalf(t, nil, err, "test1")
}

func Test_barcodeService_CustomGoods(t *testing.T) {
	params := CustomGoodsBarcodeQueryParams{
		PersonalProductSkuIdList: []int{27040554473264},
	}
	s, err := temuClient.Services.BarcodeService.CustomGoods(params)
	fmt.Println(s)
	assert.Equalf(t, nil, err, "test1")
}

func Test_barcodeService_NormalGoods(t *testing.T) {
	params := NormalGoodsBarcodeQueryParams{
		ProductSkcIdList: []int{8972250969},
	}
	s, err := temuClient.Services.BarcodeService.NormalGoods(params)
	fmt.Println(s)
	assert.Equalf(t, nil, err, "test1")
}
