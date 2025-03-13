package temu

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

// 半托管订单发货信息
type semiOrderLogisticsShipmentService service

type SemiOrderLogisticsShipmentQueryParams struct {
	ParentOrderSn string `json:"parentOrderSn"` // 父订单号
	OrderSn       string `json:"orderSn"`       // 子订单号
}

func (m SemiOrderLogisticsShipmentQueryParams) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.ParentOrderSn, validation.Required.Error("父订单号不能为空")),
		validation.Field(&m.OrderSn, validation.Required.Error("子订单号不能为空")),
	)
}

// Query 订单发货信息查询接口
// https://seller.kuajingmaihuo.com/sop/view/867739977041685428#Y01V9e
func (s *semiOrderLogisticsShipmentService) Query(ctx context.Context, params SemiOrderLogisticsShipmentQueryParams) ([]entity.ShipmentInfo, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var result = struct {
		normal.Response
		Result struct {
			ShipmentInfoDTO []entity.ShipmentInfo `json:"shipmentInfoDTO"` // 发货信息列表
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(params).
		SetResult(&result).
		Post("bg.logistics.shipment.v2.get")

	if err = recheckError(resp, result.Response, err); err != nil {
		return nil, err
	}

	return result.Result.ShipmentInfoDTO, nil
}

// 订单发货通知

type SemiOrderLogisticsShipmentConfirmInformationOrder struct {
	OrderSn       string `json:"orderSn"`       // orderSn
	ParentOrderSn string `json:"parentOrderSn"` // parentOrderSn
	GoodsId       int64  `json:"goodsId"`       // goodsId
	SkuId         int64  `json:"skuId"`         // skuId
	Quantity      int    `json:"quantity"`      // 发货数量
}
type SemiOrderLogisticsShipmentConfirmInformation struct {
	CarrierId         int64                                               `json:"carrierId"`         // 物流公司 ID
	TrackingNumber    string                                              `json:"trackingNumber"`    // 运单号
	OrderSendInfoList []SemiOrderLogisticsShipmentConfirmInformationOrder `json:"orderSendInfoList"` // 发货商品信息
}

type SemiOrderLogisticsShipmentConfirmRequest struct {
	SendType        int                                            `json:"sendType"`        // 发货类型：0-单个运单发货 1-拆成多个运单发货 2-合并发货
	SendRequestList []SemiOrderLogisticsShipmentConfirmInformation `json:"sendRequestList"` // 包裹信息
}

func (m SemiOrderLogisticsShipmentConfirmRequest) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.SendType, validation.In(0, 1, 2).Error("无效的发货类型")),
		validation.Field(&m.SendRequestList, validation.Required.Error("发货商品信息不能为空")),
	)
}

// Confirm 订单发货通知接口（bg.logistics.shipment.confirm）
// https://seller.kuajingmaihuo.com/sop/view/867739977041685428#bCMdFx
func (s *semiOrderLogisticsShipmentService) Confirm(ctx context.Context, request SemiOrderLogisticsShipmentConfirmRequest) (bool, error) {
	if err := request.validate(); err != nil {
		return false, err
	}

	var result = struct {
		normal.Response
		Result struct {
			AssistantAgreementText string   `json:"assistantAgreementText"` // 智能助手协议文案，当发货时，默认发货成功即为确认同意开启智能轨迹助手
			WarningMessage         []string `json:"warningMessage"`         // 发货操作提醒字段
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(request).
		SetResult(&result).
		Post("bg.logistics.shipment.confirm")

	if err = recheckError(resp, result.Response, err); err != nil {
		return false, err
	}

	return true, nil
}
