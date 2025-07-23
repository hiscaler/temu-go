package temu

import (
	"context"
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

// goodsPriceFullAdjustmentService 全托调价服务
type goodsPriceFullAdjustmentService service

type GoodsPriceFullAdjustmentQueryParams struct {
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

func (m GoodsPriceFullAdjustmentQueryParams) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Status, validation.In(0, 1, 2, 3).Error("无效的状态")),
	)
}

// Query 分页查询全托管调价单
// https://partner.kuajingmaihuo.com/document?cataId=875198836203&docId=908751475686
func (s goodsPriceFullAdjustmentService) Query(ctx context.Context, params GoodsPriceFullAdjustmentQueryParams) (items []entity.GoodsReviewSamplePrice, err error) {
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

type GoodsPriceFullAdjustmentConfirmItem struct {
	PriceOrderSn string `json:"priceOrderSn"` // 调价单号
	Result       int    `json:"result"`       // 审核结果 1-通过
}

func (m GoodsPriceFullAdjustmentConfirmItem) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.PriceOrderSn, validation.Required.Error("调价单号不能为空")),
		validation.Field(&m.Result, validation.In(0, 1).Error("无效的审核结果")),
	)
}

type GoodsPriceFullAdjustmentConfirmRequest struct {
	AdjustList []GoodsPriceFullAdjustmentConfirmItem `json:"adjustList"` // 调价列表
}

func (m GoodsPriceFullAdjustmentConfirmRequest) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.AdjustList,
			validation.Required.Error("调价列表不能为空"),
			validation.Each(validation.By(func(value interface{}) error {
				v, ok := value.(GoodsPriceFullAdjustmentConfirmItem)
				if !ok {
					return errors.New("无效的调价数据")
				}
				return v.validate()
			})),
		),
	)
}

// Confirm 全托管批量确认调价单
// https://partner.kuajingmaihuo.com/document?cataId=875198836203&docId=908749899377
func (s goodsPriceFullAdjustmentService) Confirm(ctx context.Context, params GoodsPriceFullAdjustmentConfirmRequest) (bool, error) {
	if err := params.validate(); err != nil {
		return false, invalidInput(err)
	}

	var result = struct {
		normal.Response
		Result any `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(params).
		SetResult(&result).
		Post("bg.full.adjust.price.batch.review")
	if err = recheckError(resp, result.Response, err); err != nil {
		return false, err
	}

	return true, nil
}
