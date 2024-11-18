package temu

import (
	"context"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
	"github.com/hiscaler/temu-go/validators/is"
)

type barcodeService service

type NormalGoodsBarcodeQueryParams struct {
	normal.ParameterWithPager
	ProductSkuIdList []int64 `json:"productSkuIdList,omitempty"` // 货品 sku id 列表
	SkcExtCode       string  `json:"skcExtCode,omitempty"`       // skc 货号
	ProductSkcIdList []int64 `json:"productSkcIdList,omitempty"` // 货品 skc id 列表
	SkuExtCode       string  `json:"skuExtCode,omitempty"`       // sku 货号
	LabelCode        int     `json:"labelCode,omitempty"`        // 标签条码
}

func (m NormalGoodsBarcodeQueryParams) Validate() error {
	return nil
}

// NormalGoods 商品条码查询v2（bg.goods.labelv2.get）
// https://seller.kuajingmaihuo.com/sop/view/889973754324016047#5LRokG
func (s barcodeService) NormalGoods(ctx context.Context, params NormalGoodsBarcodeQueryParams) (items []entity.GoodsLabel, err error) {
	params.TidyPager()
	if err = params.Validate(); err != nil {
		return
	}

	var result = struct {
		normal.Response
		Result struct {
			labelCodePageResult struct {
				TotalCount int                 `json:"totalCount"` // 总数
				Data       []entity.GoodsLabel `json:"data"`       // 结果列表
			}
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(params).
		SetResult(&result).
		Post("bg.goods.labelv2.get")
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	return result.Result.labelCodePageResult.Data, nil
}

// 定制商品条码查询（bg.goods.custom.label.get）

type CustomGoodsBarcodeQueryParams struct {
	NormalGoodsBarcodeQueryParams
	PersonalProductSkuIdList []int64 `json:"personalProductSkuIdList,omitempty"` // 定制品 sku id
}

func (m CustomGoodsBarcodeQueryParams) Validate() error {
	return nil
}

// CustomGoods 定制商品条码查询（bg.goods.custom.label.get）
// https://seller.kuajingmaihuo.com/sop/view/889973754324016047#Hc5wmR
func (s barcodeService) CustomGoods(ctx context.Context, params CustomGoodsBarcodeQueryParams) (items []entity.CustomGoodsLabel, err error) {
	params.TidyPager()
	if err = params.Validate(); err != nil {
		return
	}

	var result = struct {
		normal.Response
		Result struct {
			PersonalLabelCodePageResult struct {
				TotalCount int                       `json:"totalCount"` // 总数
				Data       []entity.CustomGoodsLabel `json:"data"`       // 结果列表
			} `json:"personalLabelCodePageResult"`
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(params).
		SetResult(&result).
		Post("bg.goods.custom.label.get")
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	return result.Result.PersonalLabelCodePageResult.Data, nil
}

// 查询箱唛（bg.logistics.boxmarkinfo.get）

type BoxMarkBarcodeQueryParams struct {
	normal.Parameter
	ReturnDataKey       bool     `json:"return_data_key"`     // 是否以打印页面url返回，如果入参是，则不返回参数信息，返回dataKey，通过拼接https://openapi.kuajingmaihuo.com/tool/print?dataKey={返回的dataKey}，访问组装的url即可打印，打印的条码按照入参参数所得结果进行打印
	DeliveryOrderSnList []string `json:"deliveryOrderSnList"` // 发货单对象列表
}

func (m BoxMarkBarcodeQueryParams) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.DeliveryOrderSnList,
			validation.Required.Error("发货单对象列表不能为空。"),
			validation.Each(validation.By(is.ShipOrderNumber())),
		),
	)
}

// BoxMarkPrintUrl 箱唛打印地址
func (s barcodeService) BoxMarkPrintUrl(ctx context.Context, deliveryOrderSnList ...string) (dataKey string, err error) {
	params := BoxMarkBarcodeQueryParams{
		ReturnDataKey:       true,
		DeliveryOrderSnList: deliveryOrderSnList,
	}
	if err = params.Validate(); err != nil {
		return
	}

	var result = struct {
		normal.Response
		Result string `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(params).
		SetResult(&result).
		Post("bg.logistics.boxmarkinfo.get")
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	return fmt.Sprintf("https://openapi.kuajingmaihuo.com/tool/print?dataKey=%s", result.Result), nil
}

// BoxMark 箱唛信息
func (s barcodeService) BoxMark(ctx context.Context, deliveryOrderSnList ...string) (items []entity.BoxMarkInfo, err error) {
	params := BoxMarkBarcodeQueryParams{
		ReturnDataKey:       false,
		DeliveryOrderSnList: deliveryOrderSnList,
	}
	if err = params.Validate(); err != nil {
		return
	}

	var result = struct {
		normal.Response
		Result []entity.BoxMarkInfo `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(params).
		SetResult(&result).
		Post("bg.logistics.boxmarkinfo.get")
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	return result.Result, nil
}
