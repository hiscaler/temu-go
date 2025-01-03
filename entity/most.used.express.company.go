package entity

// MostUsedExpressCompany 常用物流
type MostUsedExpressCompany struct {
	ServicerCode        string  `json:"servicerCode"`        // 服务商编码
	ExpressCompanyId    int64   `json:"expressCompanyId"`    // 快递公司 Id
	ExpressCompanyName  string  `json:"expressCompanyName"`  // 快递公司名称
	CanSaveChargeAmount float64 `json:"canSaveChargeAmount"` // 可节省费用（单位元）
}
