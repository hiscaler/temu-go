package entity

// GoodsParentSpecification 货品父规格
type GoodsParentSpecification struct {
	ParentSpecId   int    `json:"parentSpecId"`   //  父规格 ID
	ParentSpecName string `json:"parentSpecName"` //  父规格名称
}
