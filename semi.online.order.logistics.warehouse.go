package temu

import (
	"context"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

// 卖家发货仓库服务
type semiOnlineOrderLogisticsWarehouseService service

// Query 查询卖家发货仓库基础信息接口（bg.logistics.warehouse.list.get）
// https://seller.kuajingmaihuo.com/sop/view/144659541206936016#MdjB3d
func (s semiOnlineOrderLogisticsWarehouseService) Query(ctx context.Context) ([]entity.SemiOnlineOrderLogisticsWarehouse, error) {
	var result = struct {
		normal.Response
		Result struct {
			WarehouseList []entity.SemiOnlineOrderLogisticsWarehouse `json:"warehouseList"`
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetResult(&result).
		Post("bg.logistics.warehouse.list.get")
	if err = recheckError(resp, result.Response, err); err != nil {
		return nil, err
	}

	return result.Result.WarehouseList, nil
}
