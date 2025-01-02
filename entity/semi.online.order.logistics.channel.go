package entity

// SemiOnlineOrderLogisticsChannel 物流供应商渠道
type SemiOnlineOrderLogisticsChannel struct {
	ChannelId             int    `json:"channelId"`             // 渠道 id
	ShipLogisticsType     int    `json:"shipLogisticsType"`     // 物流产品类型
	ShipCompanyId         int    `json:"shipCompanyId"`         // 物流公司 id
	ShippingCompanyName   string `json:"shippingCompanyName"`   // 物流公司名称
	EstimatedText         string `json:"estimatedText"`         // 预估参数文案，包含预估的面单价格，币种，时效等信息 如：$39.46;USD ;1-2工作日送达
	EstimatedAmount       string `json:"estimatedAmount"`       // 预估金额，如：12.12
	EstimatedCurrencyCode string `json:"estimatedCurrencyCode"` // 预估货币币种
	InfoNeeded            []any  `json:"infoNeeded"`            // 渠道提示信息，提示这个渠道下call时候的依赖项目
	SignServiceName       string `json:"signServiceName"`       // 签收服务类型
	SignServiceId         string `json:"signServiceId"`         // 签收服务 ID
}
