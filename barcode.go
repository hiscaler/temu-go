package temu

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

type barcodeService service

type NormalGoodsBarcodeQueryParams struct {
	normal.Parameter
	Page             int    `json:"pageNo"`                     // 页码
	PageSize         int    `json:"pageSize"`                   // 页面大小
	ProductSkuIdList []int  `json:"productSkuIdList,omitempty"` // 货品sku id列表
	SkcExtCode       string `json:"skcExtCode,omitempty"`       // skc货号
	ProductSkcIdList []int  `json:"productSkcIdList,omitempty"` // 货品skc id列表
	SkuExtCode       string `json:"skuExtCode,omitempty"`       // sku货号
	LabelCode        int    `json:"labelCode,omitempty"`        // 标签条码
}

func (m NormalGoodsBarcodeQueryParams) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Page, validation.Required.Error("页码不能为空。")),
		validation.Field(&m.PageSize, validation.Required.Error("页面大小不能为空。")),
	)
}

// NormalGoods 商品条码查询v2（bg.goods.labelv2.get）
// https://seller.kuajingmaihuo.com/sop/view/889973754324016047#5LRokG
func (s barcodeService) NormalGoods(params NormalGoodsBarcodeQueryParams) (items []entity.GoodsLabel, err error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 10
	} else if params.PageSize > 500 {
		params.PageSize = 500
	}
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
	resp, err := s.httpClient.R().SetBody(params).SetResult(&result).Post("bg.goods.labelv2.get")
	if err == nil {
		err = parseResponse(resp, result.Response)
	}
	if err != nil {
		return
	}

	return result.Result.labelCodePageResult.Data, nil
}

// 定制商品条码查询（bg.goods.custom.label.get）

type CustomGoodsBarcodeQueryParams struct {
	normal.Parameter
	Page                     int    `json:"pageNo"`                             // 页码
	PageSize                 int    `json:"pageSize"`                           // 页面大小
	ProductSkuIdList         []int  `json:"productSkuIdList,omitempty"`         // 货品sku id列表
	SkcExtCode               string `json:"skcExtCode,omitempty"`               // skc货号
	ProductSkcIdList         []int  `json:"productSkcIdList,omitempty"`         // 货品skc id列表
	SkuExtCode               string `json:"skuExtCode,omitempty"`               // sku货号
	LabelCode                int    `json:"labelCode,omitempty"`                // 标签条码
	PersonalProductSkuIdList []int  `json:"personalProductSkuIdList,omitempty"` // 定制品sku id
}

func (m CustomGoodsBarcodeQueryParams) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Page, validation.Required.Error("页码不能为空。")),
		validation.Field(&m.PageSize, validation.Required.Error("页面大小不能为空。")),
	)
}

// CustomGoods 定制商品条码查询（bg.goods.custom.label.get）
// https://seller.kuajingmaihuo.com/sop/view/889973754324016047#Hc5wmR
func (s barcodeService) CustomGoods(params CustomGoodsBarcodeQueryParams) (items []entity.CustomGoodsLabel, err error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 10
	} else if params.PageSize > 500 {
		params.PageSize = 500
	}
	if err = params.Validate(); err != nil {
		return
	}
	var result = struct {
		normal.Response
		Result struct {
			labelCodePageResult struct {
				TotalCount int                       `json:"totalCount"` // 总数
				Data       []entity.CustomGoodsLabel `json:"data"`       // 结果列表
			}
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().SetBody(params).SetResult(&result).Post("bg.goods.custom.label.get")
	if err == nil {
		err = parseResponse(resp, result.Response)
	}
	if err != nil {
		return
	}

	return result.Result.labelCodePageResult.Data, nil
}

// 查询箱唛（bg.logistics.boxmarkinfo.get）

type BoxMarkBarcodeQueryParams struct {
	normal.Parameter
	ReturnDataKey       bool     `json:"return_data_key"`     // 是否以打印页面url返回，如果入参是，则不返回参数信息，返回dataKey，通过拼接https://openapi.kuajingmaihuo.com/tool/print?dataKey={返回的dataKey}，访问组装的url即可打印，打印的条码按照入参参数所得结果进行打印
	DeliveryOrderSnList []string `json:"deliveryOrderSnList"` // 发货单对象列表
}

func (m BoxMarkBarcodeQueryParams) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.DeliveryOrderSnList, validation.Required.Error("发货单对象列表不能为空。")),
	)
}

func (s barcodeService) BoxMark(params BoxMarkBarcodeQueryParams) (dataKey string, err error) {
	if err = params.Validate(); err != nil {
		return
	}
	var result = struct {
		normal.Response
		Result string `json:"result"`
	}{}
	resp, err := s.httpClient.R().SetBody(params).SetResult(&result).Post("bg.logistics.boxmarkinfo.get")
	if err == nil {
		err = parseResponse(resp, result.Response)
	}
	if err != nil {
		return
	}

	return result.Result, nil
}
