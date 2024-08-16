package entity

// LogisticsCompany 发货快递公司
type LogisticsCompany struct {
	ShipId   int    `json:"shipId"`   // 快递公司ID
	ShipName string `json:"shipName"` // 快递公司名称
}
