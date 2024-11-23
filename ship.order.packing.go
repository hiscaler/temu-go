package temu

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
	"github.com/hiscaler/temu-go/validators/is"
	"gopkg.in/guregu/null.v4"
)

// 装箱发货
type shipOrderPackingService service

// ShipOrderPackingSendRequestSelfDeliveryInformation 自行配送信息
type ShipOrderPackingSendRequestSelfDeliveryInformation struct {
	DriverUid             int    `json:"driverUid,omitempty"`             // 司机uid
	DriverName            string `json:"driverName,omitempty"`            // 司机姓名
	PlateNumber           string `json:"plateNumber,omitempty"`           // 车牌号
	DeliveryContactNumber string `json:"deliveryContactNumber,omitempty"` // 电话号码
	DeliveryContactAreaNo string `json:"deliveryContactAreaNo,omitempty"` // 电话区号
	ExpressPackageNum     int    `json:"expressPackageNum,omitempty"`     // 发货总箱数
}

// ShipOrderPackingSendRequestPlatformRecommendationDeliveryInformation 平台推荐服务商配送信息
type ShipOrderPackingSendRequestPlatformRecommendationDeliveryInformation struct {
	ExpressCompanyId          int     `json:"expressCompanyId,omitempty"`          // 快递公司Id
	TmsChannelId              int     `json:"tmsChannelId,omitempty"`              // TMS快递产品类型ID
	ExpressCompanyName        string  `json:"expressCompanyName,omitempty"`        // 快递公司名称
	StandbyExpress            bool    `json:"standbyExpress"`                      // 是否是备用快递公司
	ExpressDeliverySn         string  `json:"expressDeliverySn,omitempty"`         // 快递单号
	PredictTotalPackageWeight int64   `json:"predictTotalPackageWeight,omitempty"` // 预估总包裹重量不能为空,单位克.总量必须大于等于1千克且为整千克值
	ExpectPickUpGoodsTime     int64   `json:"expectPickUpGoodsTime,omitempty"`     // 预约取货时间
	ExpressPackageNum         int     `json:"expressPackageNum,omitempty"`         // 交接给快递公司的包裹数
	MinChargeAmount           float64 `json:"minChargeAmount,omitempty"`           // 最小预估运费（单位元）
	MaxChargeAmount           float64 `json:"maxChargeAmount,omitempty"`           // 最大预估运费（单位元）
	PredictId                 int64   `json:"predictId,omitempty"`                 // 预测ID
}

// ShipOrderPackingSendRequestThirdPartyDeliveryInformation 自行委托第三方物流配送信息
type ShipOrderPackingSendRequestThirdPartyDeliveryInformation struct {
	ExpressCompanyId   int    `json:"expressCompanyId"`            // 快递公司 Id
	ExpressCompanyName string `json:"expressCompanyName"`          // 快递公司名称
	ExpressDeliverySn  string `json:"expressDeliverySn"`           // 快递单号
	ExpressPackageNum  int    `json:"expressPackageNum,omitempty"` // 发货总箱数
}

type ShipOrderPackingSendRequest struct {
	normal.Parameter
	DeliverMethod                   null.Int                                                              `json:"deliverMethod"`                             // 发货方式
	DeliveryAddressId               int64                                                                 `json:"deliveryAddressId"`                         // 发货地址 ID
	DeliveryOrderSnList             []string                                                              `json:"deliveryOrderSnList"`                       // 发货单号
	SelfDeliveryInfo                *ShipOrderPackingSendRequestSelfDeliveryInformation                   `json:"selfDeliveryInfo,omitempty"`                // 自送信息
	ThirdPartyDeliveryInfo          *ShipOrderPackingSendRequestPlatformRecommendationDeliveryInformation `json:"thirdPartyDeliveryInfo,omitempty"`          // 公司指定物流
	ThirdPartyExpressDeliveryInfoVO *ShipOrderPackingSendRequestThirdPartyDeliveryInformation             `json:"thirdPartyExpressDeliveryInfoVO,omitempty"` // 第三方配送
}

func (m ShipOrderPackingSendRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.DeliverMethod,
			validation.Required.Error("发货方式不能为空"),
			validation.In(entity.DeliveryMethodSelf, entity.DeliveryMethodPlatformRecommendation, entity.DeliveryMethodThirdParty).Error("无效的发货方式"),
		),
		validation.Field(&m.DeliveryAddressId, validation.Required.Error("发货地址 ID 不能为空")),
		validation.Field(&m.DeliveryOrderSnList, validation.Required.Error("发货单号不能为空")),
	)
}

// Send 装箱发货接口
// https://seller.kuajingmaihuo.com/sop/view/889973754324016047#ezXrHy
func (s shipOrderPackingService) Send(ctx context.Context, request ShipOrderPackingSendRequest) (number string, err error) {
	if err = request.Validate(); err != nil {
		return
	}

	var result = struct {
		normal.Response
		Result struct {
			CreateExpressErrorRequestList []any  `json:"createExpressErrorRequestList"` // 创建快递运单失败的请求列表
			ExpressBatchSn                string `json:"expressBatchSn"`                // 创建生成的发货批次号
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(request).
		SetResult(&result).
		Post("bg.shiporder.packing.send")
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	number = result.Result.ExpressBatchSn
	return
}

// 装箱发货校验

type ShipOrderPackingMatchRequest struct {
	normal.Parameter
	DeliveryOrderSnList []string `json:"deliveryOrderSnList"` // 发货单号
}

func (m ShipOrderPackingMatchRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.DeliveryOrderSnList,
			validation.Required.Error("发货单号列表不能为空"),
			validation.Each(validation.By(is.ShipOrderNumber())),
		),
	)
}

// Match 装箱发货校验
// https://seller.kuajingmaihuo.com/sop/view/889973754324016047#TDP3qU
func (s shipOrderPackingService) Match(ctx context.Context, request ShipOrderPackingMatchRequest) (item entity.ShipOrderPackingMatchResult, err error) {
	if err = request.Validate(); err != nil {
		return
	}

	var result = struct {
		normal.Response
		Result entity.ShipOrderPackingMatchResult `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(request).
		SetResult(&result).
		Post("bg.shiporder.packing.match")
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	item = result.Result
	return
}
