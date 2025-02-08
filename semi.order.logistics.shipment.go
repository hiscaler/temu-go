package temu

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

// 半托管订单发货信息
type semiOrderLogisticsShipmentService service

type SemiOrderLogisticsShipmentGetRequest struct {
	ParentOrderSn string `json:"parentOrderSn"` // 父订单号
	OrderSn       string `json:"orderSn"`       // 子订单号
}

func (m SemiOrderLogisticsShipmentGetRequest) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.ParentOrderSn, validation.Required.Error("父订单号不能为空")),
		validation.Field(&m.OrderSn, validation.Required.Error("子订单号不能为空")),
	)
}

// Query 订单发货信息查询接口
// https://seller.kuajingmaihuo.com/sop/view/867739977041685428#Y01V9e
func (s *semiOrderLogisticsShipmentService) Query(ctx context.Context, params SemiOrderLogisticsShipmentGetRequest) ([]entity.ShipmentInfo, error) {
	if err := params.validate(); err != nil {
		return []entity.ShipmentInfo{}, err
	}

	var result = struct {
		normal.Response
		Result entity.ShipmentResult `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(params).
		SetResult(&result).
		Post("bg.logistics.shipment.get")

	if err = recheckError(resp, result.Response, err); err != nil {
		return []entity.ShipmentInfo{}, err
	}

	return result.Result.Result.ShipmentInfoDTO, nil
}
