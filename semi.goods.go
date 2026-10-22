package temu

import (
	"context"

	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

type semiGoodsService struct {
	service
	Image semiGoodsImageService
}

type SemiGoodsQueryParams struct {
	normal.ParameterWithPager
	OrderField                string  `url:"orderField,omitempty"`
	OrderType                 int     `url:"orderType,omitempty"`
	GoodsSearchType           int     `url:"goodsSearchType,omitempty"`
	SearchText                string  `url:"searchText,omitempty"`
	StatusFilterType          int     `url:"statusFilterType,omitempty"`
	CrtFrom                   int64   `url:"crtFrom,omitempty"`
	CrtTo                     int64   `url:"crtTo,omitempty"`
	GoodsIdList               []int64 `url:"goodsIdList,omitempty"`
	CatIdList                 []int64 `url:"catIdList,omitempty"`
	GoodsStatusFilterType     int     `json:"goodsStatusFilterType"`
	GoodsSubStatusFilterType  int     `json:"goodsSubStatusFilterType"`
	GoodsStatusChangeTimeFrom int64   `json:"goodsStatusChangeTimeFrom"`
	GoodsStatusChangeTimeTo   int64   `json:"goodsStatusChangeTimeTo"`
	GoodsSearchTags           []int   `json:"goodsSearchTags"`
}

func (m SemiGoodsQueryParams) validate() error {
	return nil
}

// Query 商品查询
// https://partner-us.temu.com/documentation?menu_code=fb16b05f7a904765aac4af3a24b87d4a&sub_menu_code=d2a836cf1711473ba1f83597a1b52fb0
func (s semiGoodsService) Query(ctx context.Context, params SemiGoodsQueryParams) (items []entity.SemiGoods, total, totalPages int, isLastPage bool, err error) {
	params.Page = params.TidyPager().Page
	params.OmitPage()
	if err = params.validate(); err != nil {
		return items, 0, 0, true, invalidInput(err)
	}

	var result = struct {
		normal.Response
		Result struct {
			PageNo    int                `json:"pageNo"`
			Total     int                `json:"total"`
			GoodsList []entity.SemiGoods `json:"goodsList"`
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(params).
		SetResult(&result).
		Post("bg.local.goods.list.query")
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	items = result.Result.GoodsList
	total, totalPages, isLastPage = parseResponseTotal(params.Page, params.PageSize, result.Result.Total)
	return
}

type SemiGoodsCreateRequest struct {
	Language   string     `json:"language"`
	GoodsBasic struct{}   `json:"goodsBasic"`
	Attributes []struct{} `json:"attributes"`
	SkuList    []struct{} `json:"skuList"`
}

type SemiGoodsCreateGoodsBasic struct {
}

func (m SemiGoodsCreateRequest) validate() error {
	return nil
}

// Create 添加商品
// https://partner-us.temu.com/documentation?menu_code=fb16b05f7a904765aac4af3a24b87d4a&sub_menu_code=f6f78d3f42ca41e59420cc4b7e3be479
func (s semiGoodsService) Create(ctx context.Context, req SemiGoodsCreateRequest) (int64, string, error) {
	if err := req.validate(); err != nil {
		return 0, "", invalidInput(err)
	}

	var result = struct {
		normal.Response
		Result struct {
			GoodsId         int64  `json:"goodsId"`
			ExternalGoodsId string `json:"externalGoodsId"`
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(req).
		SetResult(&result).
		Post("temu.local.goods.v3.add")
	if err = recheckError(resp, result.Response, err); err != nil {
		return 0, "", err
	}

	return result.Result.GoodsId, result.Result.ExternalGoodsId, nil
}
