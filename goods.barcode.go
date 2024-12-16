package temu

import (
	"context"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/helpers"
	"github.com/hiscaler/temu-go/normal"
	"github.com/hiscaler/temu-go/validators/is"
	"gopkg.in/guregu/null.v4"
	"time"
)

type goodsBarcodeService service

type NormalGoodsBarcodeQueryParams struct {
	normal.ParameterWithPager
	ProductSkuIdList []int64 `json:"productSkuIdList,omitempty"` // 货品 sku id 列表
	SkcExtCode       string  `json:"skcExtCode,omitempty"`       // skc 货号
	ProductSkcIdList []int64 `json:"productSkcIdList,omitempty"` // 货品 skc id 列表
	SkuExtCode       string  `json:"skuExtCode,omitempty"`       // sku 货号
	LabelCode        int     `json:"labelCode,omitempty"`        // 标签条码
}

func (m NormalGoodsBarcodeQueryParams) validate() error {
	return nil
}

// NormalGoods 商品条码查询v2（bg.goods.labelv2.get）
// https://seller.kuajingmaihuo.com/sop/view/889973754324016047#5LRokG
func (s goodsBarcodeService) NormalGoods(ctx context.Context, params NormalGoodsBarcodeQueryParams) (items []entity.GoodsLabel, err error) {
	params.TidyPager()
	if err = params.validate(); err != nil {
		err = invalidInput(err)
		return
	}

	var result = struct {
		normal.Response
		Result struct {
			LabelCodePageResult struct {
				TotalCount int                 `json:"totalCount"` // 总数
				Data       []entity.GoodsLabel `json:"data"`       // 结果列表
			} `json:"labelCodePageResult"`
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

	return result.Result.LabelCodePageResult.Data, nil
}

// 定制商品条码查询（bg.goods.custom.label.get）

type CustomGoodsBarcodeQueryParams struct {
	NormalGoodsBarcodeQueryParams
	ProductSkuIdList         []int64 `json:"productSkuIdList,omitempty"`         // 货品 SKU ID 列表
	SkcExtCode               string  `json:"skcExtCode,omitempty"`               // SKC 货号
	ProductSkcIdList         []int64 `json:"productSkcIdList,omitempty"`         // 货品 SKC ID 列表
	SkuExtCode               string  `json:"skuExtCode,omitempty"`               // SKU 货号
	LabelCode                int64   `json:"labelCode,omitempty"`                // 标签条码
	PersonalProductSkuIdList []int64 `json:"personalProductSkuIdList,omitempty"` // 定制品 SKU ID
	CreateTimeStart          string  `json:"createTimeStart,omitempty"`          // 定制品创建时间，支持毫秒时间戳
	CreateTimeEnd            string  `json:"createTimeEnd,omitempty"`            // 定制品创建时间，支持毫秒时间戳
}

func (m CustomGoodsBarcodeQueryParams) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.CreateTimeStart,
			validation.When(m.CreateTimeStart != "" || m.CreateTimeEnd != "", validation.By(is.TimeRange(m.CreateTimeStart, m.CreateTimeEnd, time.DateOnly))),
		),
	)
}

// CustomGoods 定制商品条码查询（bg.goods.custom.label.get）
// https://seller.kuajingmaihuo.com/sop/view/889973754324016047#Hc5wmR
func (s goodsBarcodeService) CustomGoods(ctx context.Context, params CustomGoodsBarcodeQueryParams) (items []entity.CustomGoodsLabel, err error) {
	params.TidyPager()
	if err = params.validate(); err != nil {
		err = invalidInput(err)
		return
	}

	if params.CreateTimeStart != "" && params.CreateTimeEnd != "" {
		if start, end, e := helpers.StrTime2UnixMilli(params.CreateTimeStart, params.CreateTimeEnd); e == nil {
			params.CreateTimeStart = start
			params.CreateTimeEnd = end
		}
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
	ReturnDataKey       null.Bool `json:"return_data_key"`     // 是否以打印页面url返回，如果入参是，则不返回参数信息，返回dataKey，通过拼接https://openapi.kuajingmaihuo.com/tool/print?dataKey={返回的dataKey}，访问组装的url即可打印，打印的条码按照入参参数所得结果进行打印
	DeliveryOrderSnList []string  `json:"deliveryOrderSnList"` // 发货单对象列表
}

func (m BoxMarkBarcodeQueryParams) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.DeliveryOrderSnList,
			validation.Required.Error("发货单对象列表不能为空"),
			validation.Each(validation.By(is.ShipOrderNumber())),
		),
	)
}

// BoxMarkPrintUrl 箱唛打印地址
func (s goodsBarcodeService) BoxMarkPrintUrl(ctx context.Context, shipOrderNumbers ...string) (dataKey string, err error) {
	params := BoxMarkBarcodeQueryParams{
		ReturnDataKey:       null.BoolFrom(true),
		DeliveryOrderSnList: shipOrderNumbers,
	}
	if err = params.validate(); err != nil {
		err = invalidInput(err)
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
func (s goodsBarcodeService) BoxMark(ctx context.Context, shipOrderNumbers ...string) (items []entity.BoxMarkInfo, err error) {
	params := BoxMarkBarcodeQueryParams{
		ReturnDataKey:       null.BoolFrom(false),
		DeliveryOrderSnList: shipOrderNumbers,
	}
	if err = params.validate(); err != nil {
		err = invalidInput(err)
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
