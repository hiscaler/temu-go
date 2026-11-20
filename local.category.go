package temu

import (
	"context"

	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

// 本本服务
type localService struct {
	Category localCategoryService // 分类
	Goods    localGoodsService    // 商品
}

type localCategoryService service

type LocalCategoryQueryParams struct {
	Language    string `json:"language,omitempty" url:"language,omitempty"`
	ParentCatId int64  `json:"parentCatId,omitempty" url:"parentCatId,omitempty"`
}

// Query 分类查询
// https://partner-us.temu.com/documentation?menu_code=fb16b05f7a904765aac4af3a24b87d4a&sub_menu_code=7ce116fe6b87443ba2a5320b25bf2b20
func (s localCategoryService) Query(ctx context.Context, params LocalCategoryQueryParams) ([]entity.LocalCategory, error) {
	var result = struct {
		normal.Response
		Result struct {
			Items []entity.LocalCategory `json:"items"`
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(params).
		SetResult(&result).
		Post("bg.local.goods.cats.get")
	if err = recheckError(resp, result.Response, err); err != nil {
		return nil, err
	}

	return result.Result.Items, nil
}
