package temu

import (
	"context"

	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

// 商品生命周期服务
type goodsLifeCycleService service

type GoodsLifeCycleQueryParams struct {
	normal.ParameterWithPager
	Page             int     `json:"pageNum"`                    // 页码
	ProductSkuIdList []int64 `json:"productSkuIdList,omitempty"` // 货品 skuId 列表
	MallId           int64   `json:"mallId,omitempty"`           // 商家店铺 ID
}

func (m GoodsLifeCycleQueryParams) validate() error {
	return nil
}

// Query 查询货品生命周期状态（bg.product.search）
// https://seller.kuajingmaihuo.com/sop/view/750197804480663142#CK9soN
func (s goodsLifeCycleService) Query(ctx context.Context, params GoodsLifeCycleQueryParams) (items []entity.GoodsLifeCycle, total, totalPages int, isLastPage bool, err error) {
	pager := params.TidyPager()
	params.Page = pager.Page
	if err = params.validate(); err != nil {
		err = invalidInput(err)
		return
	}

	var result = struct {
		normal.Response
		Result struct {
			Total    int                     `json:"total"`
			DataList []entity.GoodsLifeCycle `json:"dataList"`
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(params).
		SetResult(&result).
		Post("bg.product.search")
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	items = result.Result.DataList
	total, totalPages, isLastPage = parseResponseTotal(params.Page, params.PageSize, result.Result.Total)
	return
}
