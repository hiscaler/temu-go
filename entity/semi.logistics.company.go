package entity

// SemiLogisticsCompany 半托管物流商
type SemiLogisticsCompany struct {
	LogisticsServiceProviderId   int64  `json:"logisticsServiceProviderId"`   // 服务商 ID
	LogisticsServiceProviderName string `json:"logisticsServiceProviderName"` // 服务商名称
	LogisticsBrandName           string `json:"logisticsBrandName"`           // 服务品牌名称
}
