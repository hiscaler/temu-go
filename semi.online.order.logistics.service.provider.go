package temu

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
	"gopkg.in/guregu/null.v4"
)

// 物流服务提供商服务
type semiOnlineOrderLogisticsServiceProviderService service

type SemiOnlineOrderLogisticsServiceProviderQueryParams struct {
	WarehouseId         string     `json:"warehouseId"`                   // 仓库 id
	OrderSnList         []string   `json:"orderSnList"`                   // O 单（orderSn 非 parentOrderSn）列表（至少包含一个 O 单号）
	Weight              float64    `json:"weight"`                        // 重量（默认 2 位小数，美国lb，其他国家kg）
	WeightUnit          string     `json:"weightUnit"`                    // 重量单位
	Length              null.Float `json:"length,omitempty"`              // 包裹长度（默认 2 位小数）
	Width               float64    `json:"width"`                         // 包裹宽度（默认 2 位小数）
	Height              float64    `json:"height"`                        // 包裹高度（默认 2 位小数）
	DimensionUnit       string     `json:"dimensionUnit"`                 // 尺寸单位（美国 in，其他国家 cm）
	SignatureOnDelivery null.Bool  `json:"signatureOnDelivery,omitempty"` // 是否需要签名签收服务
}

func (m SemiOnlineOrderLogisticsServiceProviderQueryParams) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.WarehouseId, validation.Required.Error("仓库不能为空")),
		validation.Field(&m.OrderSnList, validation.Required.Error("子订单号列表不能为空")),
		validation.Field(&m.Weight,
			validation.Required.Error("重量不能为空"),
			validation.Min(0.01).Error("重量不能小于 {min}"),
		),
		validation.Field(&m.WeightUnit,
			validation.Required.Error("重量单位不能为空"),
			validation.In("lb", "kg").Error("无效的重量单位，有效值为美国 lb，其他国家 kg"),
		),
		validation.Field(&m.Height,
			validation.Required.Error("包裹长度不能为空"),
			validation.Min(0.01).Error("包裹长度不能小于 {min}"),
		),
		validation.Field(&m.Width,
			validation.Required.Error("包裹宽度不能为空"),
			validation.Min(0.01).Error("包裹宽度不能小于 {min}"),
		),
		validation.Field(&m.Height,
			validation.Required.Error("包裹高度不能为空"),

			validation.Min(0.01).Error("包裹高度不能小于 {min}"),
		),
		validation.Field(&m.DimensionUnit,
			validation.Required.Error("体积尺寸单位不能为空"),
			validation.In("in", "cm").Error("无效的体积尺寸单位，有效值为美国 in，其他国家 cm"),
		),
	)
}

// Query 查询可用物流服务接口（bg.logistics.shippingservices.get）
// https://seller.kuajingmaihuo.com/sop/view/144659541206936016#d0sexY
func (s semiOnlineOrderLogisticsServiceProviderService) Query(ctx context.Context, params SemiOnlineOrderLogisticsServiceProviderQueryParams) ([]entity.SemiOnlineOrderLogisticsChannel, error) {
	if err := params.validate(); err != nil {
		return nil, invalidInput(err)
	}

	var result = struct {
		normal.Response
		Result struct {
			OnlineChannelDtoList []entity.SemiOnlineOrderLogisticsChannel `json:"onlineChannelDtoList"` // 可使用的渠道列表
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(params).
		SetResult(&result).
		Post("bg.logistics.shippingservices.get")
	if err = recheckError(resp, result.Response, err); err != nil {
		return nil, err
	}

	items := result.Result.OnlineChannelDtoList
	for i, item := range items {
		v, e := item.ParseEstimatedAmount()
		if e != nil {
			item.AmountError = e.Error()
		} else {
			item.Amount = v
		}

		d1, d2, e := item.DeliveryDays()
		if e != nil {
			item.DeliveryDaysError = e.Error()
		} else {
			item.DeliveryMinDays = d1
			item.DeliveryMaxDays = d2
		}
		items[i] = item
	}
	return items, nil
}
