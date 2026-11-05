package entity

// SemiOrderTracking 半托管包裹物流轨迹
type SemiOrderTracking struct { // 来源：Temu temu.track.trackinginfo.get
	PackageSn    string                   `json:"packageSn"`    // 包裹号
	TrackingNum  string                   `json:"trackingNum"`  // 物流单号
	TrackingInfo []SemiOrderTrackingEvent `json:"trackingInfo"` // 物流轨迹节点列表
}

// SemiOrderTrackingEvent 物流轨迹节点
type SemiOrderTrackingEvent struct {
	LogisticsStatus    string `json:"logisticsStatus"`    // 物流状态
	LogisticsUpdatedAt string `json:"logisticsUpdatedAt"` // 物流状态更新时间
	StatusText         string `json:"statusText"`         // 物流状态文本
}
