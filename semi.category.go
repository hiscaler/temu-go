package temu

import (
	"context"

	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

// 半托管分类服务
type semiCategoryService service

type SemiCategoryQueryParams struct {
	Language    string `url:"language,omitempty"`
	ParentCatId int64  `url:"parentCatId,omitempty"`
}

// All 分类查询
// https://partner-us.temu.com/documentation?menu_code=fb16b05f7a904765aac4af3a24b87d4a&sub_menu_code=7ce116fe6b87443ba2a5320b25bf2b20
func (s semiCategoryService) Query(ctx context.Context, params SemiCategoryQueryParams) ([]entity.SemiCategory, error) {
	var result = struct {
		normal.Response
		Result struct {
			Items []entity.SemiCategory `json:"items"`
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
