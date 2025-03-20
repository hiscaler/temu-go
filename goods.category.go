package temu

import (
	"context"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
	"gopkg.in/guregu/null.v4"
)

// 商品分类服务
type goodsCategoryService struct {
	service
	Attribute goodsCategoryAttributeService
}

type GoodsCategoryQueryParams struct {
	SiteId      int       `json:"siteId"`                // 站点 id，当前固定传 1
	ParentCatId null.Int  `json:"parentCatId,omitempty"` // 上级分类 ID
	ShowHidden  null.Bool `json:"showHidden"`            // 是否展示隐藏类目，默认不展示
}

func (m GoodsCategoryQueryParams) validate() error {
	return nil
}

// Query 商品分类查询
// https://seller.kuajingmaihuo.com/sop/view/728777295758127187#G5hCmy
func (s goodsCategoryService) Query(ctx context.Context, params GoodsCategoryQueryParams) (categories []entity.Category, err error) {
	if err = params.validate(); err != nil {
		err = invalidInput(err)
		return
	}
	params.SiteId = 1 // 当前固定传 1

	var result = struct {
		normal.Response
		Result struct {
			CategoryDTOList []entity.Category `json:"categoryDTOList"`
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(params).
		SetResult(&result).
		Post("bg.goods.cats.get")
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	categories = result.Result.CategoryDTOList
	return
}
