package temu

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

// 商品数据服务
type goodsService service

type GoodsQueryParams struct {
	normal.ParameterWithPager
	SkcExtCode    string `json:"skcExtCode,omitempty"`    // 货品skc外部编码
	ProductSkcIds []int  `json:"productSkcIds,omitempty"` // SKC 列表
}

func (m GoodsQueryParams) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Page, validation.Required.Error("页码不能为空。")),
		validation.Field(&m.PageSize, validation.Required.Error("页数不能为空。")),
	)
}

// All 货品列表查询
// https://seller.kuajingmaihuo.com/sop/view/750197804480663142#SjadVR
func (s goodsService) All(ctx context.Context, params GoodsQueryParams) (items []entity.Goods, err error) {
	params.TidyPager()
	if err = params.Validate(); err != nil {
		return
	}
	var result = struct {
		normal.Response
		Result struct {
			Data       []entity.Goods `json:"data"`
			TotalCount int            `json:"totalCount"`
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(params).
		SetResult(&result).
		Post("bg.goods.list.get")
	if err == nil {
		err = parseResponse(resp, result.Response)
	}
	if err != nil {
		return
	}

	return result.Result.Data, nil
}

// One 根据商品 SKC ID 查询
func (s goodsService) One(ctx context.Context, productSkcId int) (item entity.Goods, err error) {
	items, err := s.All(ctx, GoodsQueryParams{ProductSkcIds: []int{productSkcId}})
	if err != nil {
		return
	}

	for _, v := range items {
		if v.ProductSkcId == productSkcId {
			return v, nil
		}
	}

	return item, ErrNotFound
}
