package temu

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var deliveryOrderSnList = []string{"FH2408151907154", "FH2408161777912"}

func Test_barcodeService_BoxMarkPrintUrl(t *testing.T) {
	url, err := temuClient.Services.BarcodeService.BoxMarkPrintUrl(deliveryOrderSnList...)
	assert.Equal(t, nil, err, "Services.BarcodeService.BoxMarkPrintUrl(%#v)", deliveryOrderSnList)
	assert.Contains(t, url, "https://openapi.kuajingmaihuo.com/tool/print?dataKey=")
}

func Test_barcodeService_BoxMark(t *testing.T) {
	_, err := temuClient.Services.BarcodeService.BoxMark(deliveryOrderSnList...)
	assert.Equal(t, nil, err, "Services.BarcodeService.BoxMark(%#v)", deliveryOrderSnList)
}

func Test_barcodeService_NormalGoods(t *testing.T) {
	params := NormalGoodsBarcodeQueryParams{
		ProductSkcIdList: []int{8972250969},
	}
	s, err := temuClient.Services.BarcodeService.NormalGoods(params)
	fmt.Println(s)
	assert.Equalf(t, nil, err, "test1")
}
