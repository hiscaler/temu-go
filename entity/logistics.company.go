package entity

// LogisticsCompany 发货快递公司
type LogisticsCompany struct {
	ShipId   int    `json:"shipId"`   // 快递公司 ID
	ShipName string `json:"shipName"` // 快递公司名称
}

// LogisticsExpressCompany 物流快递公司
type LogisticsExpressCompany struct {
	ExpressCompanyId   int    `json:"expressCompanyId"`   // 快递公司 Id
	ExpressCompanyName string `json:"expressCompanyName"` // 快递公司名称
}
