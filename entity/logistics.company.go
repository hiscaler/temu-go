package entity

// LogisticsShippingCompany 发货快递公司
type LogisticsShippingCompany struct {
	ShipId   int64  `json:"shipId"`   // 快递公司 ID
	ShipName string `json:"shipName"` // 快递公司名称
}

// LogisticsExpressCompany 物流快递公司
type LogisticsExpressCompany struct {
	ExpressCompanyId   int64  `json:"expressCompanyId"`   // 快递公司 ID
	ExpressCompanyName string `json:"expressCompanyName"` // 快递公司名称
}
