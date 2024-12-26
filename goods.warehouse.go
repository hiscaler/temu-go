package temu

import (
	"context"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

// 商品仓库服务
type goodsWarehouseService service

type GoodsWarehouseQueryParams struct {
	SiteIdList []int64 `json:"siteIdList"` // 站点列表
}

func (m GoodsWarehouseQueryParams) validate() error {
	return nil
}

// Query 根据站点查询可绑定的发货仓库信息接口
// https://seller.kuajingmaihuo.com/sop/view/867739977041685428#hpIdAo
func (s goodsWarehouseService) Query(ctx context.Context, params GoodsWarehouseQueryParams) (items []entity.SiteWarehouses, err error) {
	if err = params.validate(); err != nil {
		err = invalidInput(err)
		return
	}
	var result = struct {
		normal.Response
		Result struct {
			WarehouseDTOList []entity.SiteWarehouses `json:"warehouseDTOList"` // 可选发货仓列表
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(params).
		SetResult(&result).
		Post("bg.goods.warehouse.list.get")
	if err = recheckError(resp, result.Response, err); err != nil {
		return items, err
	}

	return result.Result.WarehouseDTOList, nil
}
