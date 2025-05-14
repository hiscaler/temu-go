package temu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var deliveryOrderSnList = []string{"FH2411282498847"}

func Test_goodsBarcodeService_BoxMarkPrintUrl(t *testing.T) {
	url, err := temuClient.Services.Goods.Barcode.BoxMarkPrintUrl(ctx, deliveryOrderSnList...)
	assert.Equalf(t, nil, err, "Services.Goods.Barcode.BoxMarkPrintUrl(ctx, %#v)", deliveryOrderSnList)
	assert.Contains(t, url, "https://openapi.kuajingmaihuo.com/tool/print?dataKey=")
}

func Test_goodsBarcodeService_BoxMark(t *testing.T) {
	_, err := temuClient.Services.Goods.Barcode.BoxMark(ctx, deliveryOrderSnList...)
	assert.Equalf(t, nil, err, "Services.Goods.Barcode.BoxMark(ctx, %#v)", deliveryOrderSnList)
}

func Test_goodsBarcodeService_NormalGoods(t *testing.T) {
	params := NormalGoodsBarcodeQueryParams{
		ProductSkcIdList: []int64{92543270006},
	}
	_, err := temuClient.Services.Goods.Barcode.NormalGoods(ctx, params)
	assert.Equalf(t, nil, err, "Services.Goods.Barcode.NormalGoods(ctx, %#v)", params)
}

func Test_goodsBarcodeService_NormalGoodsPrintUrl(t *testing.T) {
	params := NormalGoodsBarcodeQueryParams{
		ProductSkcIdList: []int64{92543270006},
	}

	url, err := temuClient.Services.Goods.Barcode.NormalGoodsPrintUrl(ctx, params)
	assert.Equalf(t, nil, err, "Services.Goods.Barcode.NormalGoodsPrintUrl(ctx, %#v)", params)
	assert.Contains(t, url, "https://openapi.kuajingmaihuo.com/tool/print?dataKey=")
}

func Test_goodsBarcodeService_CustomGoods(t *testing.T) {
	params := CustomGoodsBarcodeQueryParams{
		PersonalProductSkuIdList: []int64{60294097402138},
	}
	items, err := temuClient.Services.Goods.Barcode.CustomGoods(ctx, params)
	assert.Equalf(t, nil, err, "Services.Goods.Barcode.CustomGoods(ctx, %#v)", params)
	assert.Equalf(t, true, len(items) != 0, "Services.Goods.Barcode.CustomGoods(ctx, %#v)", params)
}
