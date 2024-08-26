package temu

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

// 物流服务

type logisticsService service

// Companies 查询发货快递公司
// https://seller.kuajingmaihuo.com/sop/view/889973754324016047#wjtGTK
func (s logisticsService) Companies(ctx context.Context) (items []entity.LogisticsCompany, err error) {
	var result = struct {
		normal.Response
		Result struct {
			ShipList []entity.LogisticsCompany `json:"shipList"` // 快递公司列表
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetResult(&result).
		Post("bg.logistics.company.get")
	if err == nil {
		err = parseResponse(resp, result.Response)
	}
	if err != nil {
		return
	}

	return result.Result.ShipList, nil
}

// Company 根据 ID 查询发货快递公司
func (s logisticsService) Company(ctx context.Context, shipId int) (item entity.LogisticsCompany, err error) {
	items, err := s.Companies(ctx)
	if err != nil {
		return
	}

	for _, company := range items {
		if company.ShipId == shipId {
			return company, nil
		}
	}
	return item, ErrNotFound
}

// 平台推荐物流商匹配接口
// https://seller.kuajingmaihuo.com/sop/view/889973754324016047#16WiXI

type LogisticsMatchRequest struct {
	DeliveryAddressId         int                   `json:"deliveryAddressId"`         // 发货地址
	PredictTotalPackageWeight int                   `json:"predictTotalPackageWeight"` // 预估总包裹重量，单位g
	UrgencyType               int                   `json:"urgencyType"`               // 是否是紧急发货单，0-普通 1-急采
	SubWarehouseId            int                   `json:"subWarehouseId"`            // 收货子仓id
	QueryStandbyExpress       bool                  `json:"queryStandbyExpress"`       // 是否查询备用快递服务商, false-不查询 true-查询
	TotalPackageNum           int                   `json:"totalPackageNum"`           // 包裹件数
	ReceiveAddressInfo        entity.ReceiveAddress `json:"receiveAddressInfo"`        // 收货地址
	DeliveryOrderSns          []string              `json:"deliveryOrderSns"`
}

func (m LogisticsMatchRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.PredictTotalPackageWeight, validation.Required.Error("预估总包裹重量不能为空。")),
		validation.Field(&m.TotalPackageNum, validation.Required.Error("包裹件数不能为空。")),
	)
}

func (s logisticsService) Match(ctx context.Context, request LogisticsMatchRequest) (items []entity.LogisticsMatch, err error) {
	if err = request.Validate(); err != nil {
		return
	}

	var result = struct {
		normal.Response
		Result []entity.LogisticsMatch `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(request).
		SetResult(&result).
		Post("bg.shiporderv2.logisticsmatch.get")
	if err == nil {
		err = parseResponse(resp, result.Response)
	}
	if err != nil {
		return
	}

	return result.Result, nil
}
