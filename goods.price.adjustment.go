package temu

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

// goodsPriceAdjustmentService 调价服务
type goodsPriceAdjustmentService service

// 分页查询全托管调价单
// https://partner.kuajingmaihuo.com/document?cataId=875198836203&docId=908751475686

type GoodsPriceAdjustmentQueryParams struct {
	normal.ParameterWithPager
	SkcId               []int64  `json:"skcId"`               // Skc ID
	ExtCodeType         int      `json:"extCodeType"`         // 货号类型 1-skc货号 2-sku货号
	PriceOrderSn        []string `json:"priceOrderSn"`        // 调价单号
	PriceType           int      `json:"priceType"`           // 价格类型, 0-日常价，1-活动价
	FilterProductSource int      `json:"filterProductSource"` // 调价来源. 可选值含义说明:[1:超越爆款计划;2:绿通优先发货;3:其他;47:广告渠道补贴;]
	Source              int      `json:"source"`              // 申请来源 1-运营，2-供应商
	ProductName         string   `json:"productName"`         // 货品名称
	SupportPersonal     int      `json:"supportPersonal"`     // 是否支持定制化商品，1的时候是定制,0 是查非定制，为空不做筛选
	ExtCodes            []string `json:"extCodes"`            // 货号列表
	CreatedAtEnd        int64    `json:"createdAtEnd"`        // 创建日期-结束，精确到毫秒(13位)
	CreatedAtBegin      int64    `json:"createdAtBegin"`      // 创建日期-开始，精确到毫秒(13位)
	Status              int      `json:"status"`              // 状态 0-待调价，1-带供应商确认，2-调价成功，3-调价失败
}

func (m GoodsPriceAdjustmentQueryParams) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Status, validation.In(0, 1, 2, 3).Error("无效的状态")),
	)
}

func (s goodsPriceAdjustmentService) Query(ctx context.Context, params GoodsPriceAdjustmentQueryParams) (items []entity.GoodsReviewSamplePrice, err error) {
	if err = params.validate(); err != nil {
		return items, invalidInput(err)
	}

	var result = struct {
		normal.Response
		Result struct {
			Total int                             `json:"total"`
			List  []entity.GoodsReviewSamplePrice `json:"list"`
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(params).
		SetResult(&result).
		Post("bg.full.adjust.price.page.query")
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	return result.Result.List, nil
}
