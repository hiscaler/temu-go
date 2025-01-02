package temu

import (
	"context"
	"github.com/hiscaler/temu-go/normal"
	"gopkg.in/guregu/null.v4"
)

// 物流发货服务
type semiOnlineOrderLogisticsShipmentService service

type SemiLogisticsShipmentCreateRequest struct {
	SendType int `json:"sendType"` // 发货类型：0-单个运单发货 1-拆成多个运单发货 2-合并发货
	// TRUE：下call成功之后延迟发货
	// FALSE/不填：下call成功订单自动流转为已发货
	ShipLater          bool   `json:"shipLater,omitempty"`          // 下 call 成功后是否延迟发货
	ShipLaterLimitTime string `json:"shipLaterLimitTime,omitempty"` // 稍后发货兜底配置时间（单位:h），枚举：24, 48, 72, 96
	SendRequestList    []struct {
		ShipCompanyId     int64  `json:"shipCompanyId"`  // 物流公司 id
		TrackingNumber    string `json:"trackingNumber"` // 运单号
		OrderSendInfoList []struct {
			OrderSn       string `json:"orderSn"`       // 订单号
			ParentOrderSn string `json:"parentOrderSn"` // 父订单号
			GoodsId       string `json:"goodsId"`       // 商品 goodsId
			SkuId         int64  `json:"skuId"`         // 商品 skuId
			Quantity      int    `json:"quantity"`      // 发货数量
		} `json:"orderSendInfoList"` // 发货商品信息
		WarehouseId        int64  `json:"warehouseId"`     // 仓库id
		Weight             string `json:"weight"`          // 重量（默认2位小数）
		WeightUnit         string `json:"weightUnit"`      // 重量单位，美国为lb（磅），其他国家为kg（千克）
		Length             string `json:"length"`          // 包裹长度（默认2位小数）
		Width              string `json:"width"`           // 包裹宽度（默认2位小数）
		Height             string `json:"height"`          // 包裹高度（默认2位小数）
		DimensionUnit      string `json:"dimensionUnit"`   // 尺寸单位高度 ，美国为in（英寸）其他国家为cm（厘米）
		ChannelId          int    `json:"channelId"`       // 渠道id，取自shipservice.get
		PickupStartTime    int64  `json:"pickupStartTime"` // 预约上门取件开始时间（当渠道为需要下call同时入参预约时间渠道时，需入参。剩余渠道无需入参。）
		PickupEndTime      int64  `json:"pickupEndTime"`   // 预约上门取件结束时间（当渠道为需要下call同时入参预约时间渠道时，需入参。剩余渠道无需入参。）
		SignServiceId      int64  `json:"signServiceId"`   // 想使用的签收服务ID
		SplitSubPackage    bool   `json:"splitSubPackage"` // 是否为单件SKU拆多包裹（TRUE：是单件SKU多包裹场景 FALSE/不填：不是单件SKU多包裹场景）
		SendSubRequestList []struct {
			ExtendWeightUnit string `json:"extendWeightUnit"` // 扩展重量单位
			ExtendWeight     string `json:"extendWeight"`     // 扩展重量
			WeightUnit       string `json:"weightUnit"`       // 重量单位
			DimensionUnit    string `json:"dimensionUnit"`    // 尺寸单位
			Length           string `json:"length"`           // 包裹长度（默认2位小数）
			Weight           string `json:"weight"`           // 包裹重量（默认2位小数）
			Height           string `json:"height"`           // 包裹宽度（默认2位小数）
			WarehouseId      string `json:"warehouseId"`      // 仓库id
			ShipCompanyId    string `json:"shipCompanyId"`    // 物流公司ID
			SignServiceId    int64  `json:"signServiceId"`    // 想使用的签收服务ID
		} `json:"sendSubRequestList"` // 单件sku多包裹场景，附属包裹入参
	} `json:"sendRequestList"` // 包裹信息
}

func (m SemiLogisticsShipmentCreateRequest) validate() error {
	return nil
}

// Create 物流在线发货下单接口（bg.logistics.shipment.create）
// https://seller.kuajingmaihuo.com/sop/view/144659541206936016#Tf6UNY
func (s semiOnlineOrderLogisticsShipmentService) Create(ctx context.Context, request SemiLogisticsShipmentCreateRequest) (items []string, limitTime null.String, err error) {
	if err = request.validate(); err != nil {
		return
	}

	var result = struct {
		normal.Response
		Result struct {
			PackageSnList      []string    `json:"packageSnList"`      // 可使用的渠道列表
			ShipLaterLimitTime null.String `json:"shipLaterLimitTime"` // 稍后发货兜底配置时间，如下 call 时有则返回
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetResult(&result).
		Post("bg.logistics.shipment.create")
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	return result.Result.PackageSnList, result.Result.ShipLaterLimitTime, nil
}
