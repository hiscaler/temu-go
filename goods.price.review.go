package temu

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

// goodsPriceReviewService 供货价/核价/调价服务
type goodsPriceReviewService service

// 分页查询核价单
// https://partner.kuajingmaihuo.com/document?cataId=875198836203&docId=899321422992

type GoodsPriceReviewQueryParams struct {
	normal.ParameterWithPager
	IdLt            int64 `json:"idLt,omitempty"`            // id 范围查询最大值
	IdGt            int64 `json:"idGt,omitempty"`            // id 范围查询最小值
	OrderStatusList []int `json:"orderStatusList,omitempty"` // 核价单状态列表. 可选值含义说明:[0:待核价;1:待供应商确认;2:核价通过;3:核价驳回;4:废弃;5:价格同步中;]
}

func (m GoodsPriceReviewQueryParams) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.OrderStatusList, validation.In(0, 1, 2, 3, 4, 5).Error("无效的核价单状态")),
	)
}

func (s goodsPriceReviewService) Query(ctx context.Context, params GoodsPriceReviewQueryParams) (items []entity.GoodsReviewSamplePrice, err error) {
	if err = params.validate(); err != nil {
		return items, invalidInput(err)
	}

	var result = struct {
		normal.Response
		Result struct {
			Total                 int                             `json:"total"`
			ReviewSamplePriceList []entity.GoodsReviewSamplePrice `json:"reviewSamplePriceList"` // 核价单 sku 及建议价
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(params).
		SetResult(&result).
		Post("bg.price.review.page.query")
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	return result.Result.ReviewSamplePriceList, nil
}

// Confirm 同意核价单建议价
// https://partner.kuajingmaihuo.com/document?cataId=875198836203&docId=901412462419
func (s goodsPriceReviewService) Confirm(ctx context.Context, orderId int64) (bool, error) {
	var result = struct {
		normal.Response
		Result any `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(map[string]int64{"orderId": orderId}).
		SetResult(&result).
		Post("bg.price.review.confirm")
	if err = recheckError(resp, result.Response, err); err != nil {
		return false, err
	}

	return true, nil
}

// Reject 不同意核价单建议价（并给出新的申报价）
// https://partner.kuajingmaihuo.com/document?cataId=875198836203&docId=901413494559

type GoodsPriceReviewRejectRequest struct {
	OrderId           int64 `json:"orderId"` // 核价单 ID
	BargainReasonList []struct {
		Type   int    `json:"type"`   // 重新报价原因类型. 可选值含义说明:[0:材质;1:功能;2:其他;3:品类;4:外观;5:版型;6:图案;7:规格尺寸;8:品牌;]
		Reason string `json:"reason"` // 具体原因
	} `json:"bargainReasonList,omitempty"` // 重新报价原因列表
	ExternalLinkList []string `json:"externalLinkList"` // 外部链接，最多录入5个链接
	PriceItemList    []struct {
		ProductSkuId int64   `json:"productSkuId"`
		Price        float64 `json:"price"`
	} `json:"priceItemList,omitempty"`
}

func (m GoodsPriceReviewRejectRequest) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.OrderId, validation.Required.Error("核价单 ID 不能为空")),
	)
}

func (s goodsPriceReviewService) Reject(ctx context.Context, request GoodsPriceReviewRejectRequest) (bool, error) {
	var result = struct {
		normal.Response
		Result any `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(request).
		SetResult(&result).
		Post("bg.price.review.reject")
	if err = recheckError(resp, result.Response, err); err != nil {
		return false, err
	}

	return true, nil
}
