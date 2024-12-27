package temu

import (
	"context"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

// 发货仓库服务（半托管专属）
type semiLogisticsWarehouseService service

// Query 查询卖家发货仓库基础信息接口（bg.logistics.warehouse.list.get）
// https://seller.kuajingmaihuo.com/sop/view/144659541206936016#MdjB3d
func (s semiLogisticsWarehouseService) Query(ctx context.Context) (items []entity.SemiLogisticsWarehouse, err error) {
	var result = struct {
		normal.Response
		Result struct {
			WarehouseList []entity.SemiLogisticsWarehouse `json:"warehouseList"`
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetResult(&result).
		Post("bg.logistics.warehouse.list.get")
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	items = result.Result.WarehouseList
	return
}
