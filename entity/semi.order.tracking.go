package entity

// SemiOrderTracking 半托管包裹物流轨迹。
type SemiOrderTracking struct {
	PackageSn    string                   `json:"packageSn"`
	TrackingNum  string                   `json:"trackingNum"`
	TrackingInfo []SemiOrderTrackingEvent `json:"trackingInfo"`
}

// SemiOrderTrackingEvent 物流轨迹节点。
type SemiOrderTrackingEvent struct {
	LogisticsStatus    string `json:"logisticsStatus"`
	LogisticsUpdatedAt string `json:"logisticsUpdatedAt"`
	StatusText         string `json:"statusText"`
}
