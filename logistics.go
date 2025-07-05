package temu

import (
	"context"

	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

// 物流服务

type logisticsService service

// Companies 查询发货快递公司
// https://seller.kuajingmaihuo.com/sop/view/889973754324016047#wjtGTK
func (s logisticsService) Companies(ctx context.Context) (items []entity.LogisticsShippingCompany, err error) {
	var result = struct {
		normal.Response
		Result struct {
			ShipList []entity.LogisticsShippingCompany `json:"shipList"` // 快递公司列表
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetResult(&result).
		Post("bg.logistics.company.get")
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	return result.Result.ShipList, nil
}

// Company 根据 ID 查询发货快递公司
func (s logisticsService) Company(ctx context.Context, shippingId int64) (item entity.LogisticsShippingCompany, err error) {
	items, err := s.Companies(ctx)
	if err != nil {
		return
	}

	for _, company := range items {
		if company.ShipId == shippingId {
			return company, nil
		}
	}
	return item, ErrNotFound
}
