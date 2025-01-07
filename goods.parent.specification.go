package temu

import (
	"context"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

// 商品父规格列表服务
type goodsParentSpecificationService service

// Query 查询父规格列表
func (s goodsParentSpecificationService) Query(ctx context.Context) ([]entity.GoodsParentSpecification, error) {
	var result = struct {
		normal.Response
		Result struct {
			ParentSpecDTOS []entity.GoodsParentSpecification `json:"parentSpecDTOS"`
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetResult(&result).
		Post("bg.goods.parentspec.get")
	if err = recheckError(resp, result.Response, err); err != nil {
		return nil, err
	}

	return result.Result.ParentSpecDTOS, nil
}
