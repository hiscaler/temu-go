package temu

import (
	"context"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

// 商品数据服务
type goodsService service

type GoodsQueryParams struct {
	normal.ParameterWithPager
	Page                   int     `json:"page"`                             // 页码
	Cat1Id                 int     `json:"cat1Id,omitempty"`                 // 一级分类 ID
	Cat2Id                 int     `json:"cat2Id,omitempty"`                 // 二级分类 ID
	Cat3Id                 int     `json:"cat3Id,omitempty"`                 // 三级分类 ID
	Cat4Id                 int     `json:"cat4Id,omitempty"`                 // 四级分类 ID
	Cat5Id                 int     `json:"cat5Id,omitempty"`                 // 五级分类 ID
	Cat6Id                 int     `json:"cat6Id,omitempty"`                 // 六级分类 ID
	Cat7Id                 int     `json:"cat7Id,omitempty"`                 // 七级分类 ID
	Cat8Id                 int     `json:"cat8Id,omitempty"`                 // 八级分类 ID
	Cat9Id                 int     `json:"cat9Id,omitempty"`                 // 九级分类 ID
	Cat10Id                int     `json:"cat10Id,omitempty"`                // 十级分类 ID
	SkcExtCode             string  `json:"skcExtCode,omitempty"`             // 货品 SKC 外部编码
	SupportPersonalization int     `json:"supportPersonalization,omitempty"` // 是否支持定制品模板
	BindSiteIds            []int   `json:"bindSiteIds,omitempty"`            // 经营站点
	ProductName            string  `json:"productName,omitempty"`            // 货品名称
	ProductSkcIds          []int64 `json:"productSkcIds,omitempty"`          // SKC 列表
	CreatedAtStart         int     `json:"createdAtStart,omitempty"`         // 创建时间开始，毫秒级时间戳
	CreatedAtEnd           int     `json:"createdAtEnd,omitempty"`           // 创建时间结束，毫秒级时间戳
}

func (m GoodsQueryParams) validate() error {
	return nil
}

// Query 货品列表查询
// https://seller.kuajingmaihuo.com/sop/view/750197804480663142#SjadVR
func (s goodsService) Query(ctx context.Context, params GoodsQueryParams) (items []entity.Goods, total, totalPages int, isLastPage bool, err error) {
	params.TidyPager()
	if err = params.validate(); err != nil {
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
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	items = result.Result.Data
	total, totalPages, isLastPage = parseResponseTotal(params.Page, params.PageSize, result.Result.TotalCount)
	return
}

// One 根据商品 SKC ID 查询
func (s goodsService) One(ctx context.Context, productSkcId int64) (item entity.Goods, err error) {
	items, _, _, _, err := s.Query(ctx, GoodsQueryParams{ProductSkcIds: []int64{productSkcId}})
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
