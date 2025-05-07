package entity

import "gopkg.in/guregu/null.v4"

// SemiOnlineOrderLogisticsShipmentPackageResult 物流在线发货下单结果
type SemiOnlineOrderLogisticsShipmentPackageResult struct {
	PackageSn             string      `json:"packageSn"` // 包裹号
	EstimatedText         string      `json:"estimatedText"`
	EstimatedCurrencyCode string      `json:"estimatedCurrencyCode"`
	ExtendWeight          null.String `json:"extendWeight"`
	ExtendWeightUnit      null.String `json:"extendWeightUnit"`
	// 包裹主子类型，对应枚举如下：
	// 主包裹：MAIN
	// 子包裹：SPLIT_LARGE_ITEM
	// 备注：
	// 1、当包裹下call是单件SKU多包裹的形式时，会将【sendSubRequestList】入参的相关包裹视为关联这个包裹的子包裹。
	// 2、默认所有非单sku多包裹场景的包裹都是主包裹。
	SubPackageType      string      `json:"subPackageType"`      // 包裹主子类型
	MainPackageSn       string      `json:"mainPackageSn"`       // 该包裹关联的主包裹。 当包裹为主包裹时返回自身，当包裹为子包裹时返回关联的主包裹
	SubPackageSnList    []string    `json:"subPackageSnList"`    // 该包裹下的子包裹列表。当包裹为主包裹时返回关联的子包裹，当包裹为子包裹时返回为空。
	PackageDeliveryType int         `json:"packageDeliveryType"` // 包裹发货类型，可选值含义说明:[0:不适用（未知的）;1:商家发货;2:平台发货;
	WarehouseId         string      `json:"warehouseId"`         // 仓库id
	WarehouseName       string      `json:"warehouseName"`       // 仓库名称
	FailReasonText      null.String `json:"failReasonText"`      // 下单失败的原因 当命中周六不可派送场景时，下单失败原因提示如下 “Saturday Delivery is unavailable to desired destination.”
	// 当命中周六不可派送场景时，展示解决方案提示如下
	// "On Saturdays, home delivery is not allowed. If you confirm that you agree not to have home delivery on Saturdays, can call the bg.logistics.shipment.update and retry placing an order with Confirm acceptance."
	ReservationSn  null.String `json:"reservationSn"`
	SolutionText   null.String `json:"solutionText"`   // 失败原因对应的解决方案
	WarningMessage []string    `json:"warningMessage"` // 提醒信息
	// 对应枚举如下：0申请中 1成功 2失败
	// 当且仅当状态为2失败才能再次在线下单
	// 状态为1成功和2失败都能转为自发货
	ShippingLabelStatus   int    `json:"shippingLabelStatus"`   // 当前包裹下单状态
	CanChangeToManualSend bool   `json:"canChangeToManualSend"` // 是否可以修改成手动填写物流单号
	Weight                string `json:"weight"`                // 重量（默认 2 位小数）
	WeightUnit            string `json:"weightUnit"`            // 重量单位
	Length                string `json:"length"`                // 包裹长度（默认 2 位小数）
	Width                 string `json:"width"`                 // 包裹宽度（默认 2 位小数）
	Height                string `json:"height"`                // 包裹高度（默认 2 位小数）
	DimensionUnit         string `json:"dimensionUnit"`         // 尺寸单位高度
	ChannelId             int64  `json:"channelId"`             // 渠道 id
	ShipCompanyId         int64  `json:"shipCompanyId"`         // 物流公司 id
	ShipLogisticsType     string `json:"shipLogisticsType"`     // 物流公司类型
	ShippingCompanyName   string `json:"shippingCompanyName"`   // 运单号
	TrackingNumber        string `json:"trackingNumber"`        // 物流公司名称
	OrderSendInfoList     []struct {
		ParentOrderSn string `json:"parentOrderSn"` // PO 单号
		OrderSn       string `json:"orderSn"`       // O 单号
		GoodsId       int64  `json:"goodsId"`       // goodsId
		SkuId         int64  `json:"skuId"`         // skuId
		Quantity      int    `json:"quantity"`      // 数量
	} `json:"orderSendInfoList"` // 订单信息
	SignServiceId   null.String `json:"signServiceId"`
	PickupStartTime null.Time   `json:"pickupStartTime"`
	PickupEndTime   null.Time   `json:"pickupEndTime"`
	EstimatedAmount null.String `json:"estimatedAmount"`
}
