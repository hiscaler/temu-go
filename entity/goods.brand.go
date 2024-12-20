package entity

// GoodsBrand 货品品牌
type GoodsBrand struct {
	BrandId       int64  `json:"brandId"`       // 品牌 ID
	BrandNameEn   string `json:"brandNameEn"`   // 品牌英文名
	Pid           int64  `json:"pid"`           // 基础属性值 ID
	Vid           int64  `json:"vid"`           // 属性值 ID
	RegSerialCode string `json:"regSerialCode"` // 注册序列号
}
