package entity

// MallType 店铺类型
type MallType struct {
	IsSemiManagedMall bool `json:"semiManagedMall"` // 是否是半托管店铺
	IsThriftStore     bool `json:"isThriftStore"`   // 是否是二手店
}
