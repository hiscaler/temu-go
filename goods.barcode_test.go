package temu

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var deliveryOrderSnList = []string{"FH2411282498847"}

func Test_goodsBarcodeService_BoxMarkPrintUrl(t *testing.T) {
	url, err := temuClient.Services.Goods.Barcode.BoxMarkPrintUrl(ctx, deliveryOrderSnList...)
	assert.Equal(t, nil, err, "Services.Goods.Barcode.BoxMarkPrintUrl(ctx, %#v)", deliveryOrderSnList)
	assert.Contains(t, url, "https://openapi.kuajingmaihuo.com/tool/print?dataKey=")
}

func Test_goodsBarcodeService_BoxMark(t *testing.T) {
	_, err := temuClient.Services.Goods.Barcode.BoxMark(ctx, deliveryOrderSnList...)
	assert.Equal(t, nil, err, "Services.Goods.Barcode.BoxMark(ctx, %#v)", deliveryOrderSnList)
}

func Test_goodsBarcodeService_NormalGoods(t *testing.T) {
	params := NormalGoodsBarcodeQueryParams{
		ProductSkcIdList: []int64{8972250969},
	}
	_, err := temuClient.Services.Goods.Barcode.NormalGoods(ctx, params)
	assert.Equalf(t, nil, err, "Services.Goods.Barcode.NormalGoods(ctx, %#v)", params)
}

func Test_goodsBarcodeService_CustomGoods(t *testing.T) {
	params := CustomGoodsBarcodeQueryParams{
		PersonalProductSkuIdList: []int64{60294097402138},
	}
	items, err := temuClient.Services.Goods.Barcode.CustomGoods(ctx, params)
	assert.Equalf(t, nil, err, "Services.Goods.Barcode.CustomGoods(ctx, %#v)", params)
	assert.Equal(t, true, len(items) != 0, "Services.Goods.Barcode.CustomGoods(%#v)", params)
}
