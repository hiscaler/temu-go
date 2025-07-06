package temu

import (
	"context"

	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

// 商品品牌数据服务
type goodsBrandService service

// 货品品牌

type GoodsBrandQueryParams struct {
	normal.ParameterWithPager
	Page      int    `json:"page"`                // 页码
	Vid       int64  `json:"vid,omitempty"`       // 搜索的属性 ID
	BrandName string `json:"BrandName,omitempty"` // 搜索的品牌名称
}

func (m GoodsBrandQueryParams) validate() error {
	return nil
}

// Query 查询可绑定的品牌接口
// https://seller.kuajingmaihuo.com/sop/view/750197804480663142#PjxWnZ
func (s goodsBrandService) Query(ctx context.Context, params GoodsBrandQueryParams) (items []entity.GoodsBrand, total, totalPages int, isLastPage bool, err error) {
	params.Page = params.TidyPager().Page
	params.OmitPage()
	if err = params.validate(); err != nil {
		return items, 0, 0, true, invalidInput(err)
	}

	var result = struct {
		normal.Response
		Result struct {
			Total     int                 `json:"total"`
			PageItems []entity.GoodsBrand `json:"pageItems"`
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(params).
		SetResult(&result).
		Post("bg.goods.brand.get")
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	items = result.Result.PageItems
	total, totalPages, isLastPage = parseResponseTotal(params.Page, params.PageSize, result.Result.Total)
	return
}
