package temu

import (
	"context"
	"github.com/hiscaler/temu-go/normal"
)

// 商品仓库服务
type goodsWarehouseService service

type GoodsWarehouseQueryParams struct {
	SiteIdList []int64 `json:"siteIdList"`
}

func (m GoodsWarehouseQueryParams) validate() error {
	return nil
}

func (s goodsWarehouseService) Query(ctx context.Context, params GoodsWarehouseQueryParams) (err error) {
	if err = params.validate(); err != nil {
		err = invalidInput(err)
		return
	}
	var result = struct {
		normal.Response
		Result any `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(params).
		SetResult(&result).
		Post("bg.goods.warehouse.list.get")
	if err = recheckError(resp, result.Response, err); err != nil {
		return err
	}

	return nil
}
