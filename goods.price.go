package temu

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

// goodsPriceService 供货价/核价/调价服务
type goodsPriceService struct {
	service
	Review         goodsPriceReviewService
	FullAdjustment goodsPriceAdjustmentService
}

// Query 货品供货价查询
// https://partner.kuajingmaihuo.com/document?cataId=875198836203&docId=901410718805

type GoodsPriceQueryParams struct {
	ProductSkuIds []int64 `json:"productSkuIds"` // 货品 sku ID
}

func (m GoodsPriceQueryParams) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.ProductSkuIds, validation.Required.Error("货品 sku ID 列表不能为空")),
	)
}

func (s goodsPriceService) Query(ctx context.Context, params GoodsPriceQueryParams) (items []entity.ProductSkuSupplierPrice, err error) {
	if err = params.validate(); err != nil {
		return items, invalidInput(err)
	}

	var result = struct {
		normal.Response
		Result struct {
			ProductSkuSupplierPriceList []entity.ProductSkuSupplierPrice `json:"productSkuSupplierPriceList"`
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(params).
		SetResult(&result).
		Post("bg.goods.price.list.get")
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	return result.Result.ProductSkuSupplierPriceList, nil
}
