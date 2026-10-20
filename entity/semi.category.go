package entity

// SemiCategory 半托管类目
type SemiCategory struct {
	CatId              int    `json:"catId"`              // 类目 ID
	AvailableStatus    int    `json:"availableStatus"`    // 可用状态 1=可用，2=不可用
	Level              int    `json:"level"`              // 类目层级
	CatName            string `json:"catName"`            // 类目名称
	SecondHandCategory bool   `json:"secondHandCategory"` // 是否二手类目
	CatType            int    `json:"catType"`            // 类目类型 1=服饰，其他的非服饰
	Leaf               bool   `json:"leaf"`               // 是否叶子类目
	ParentId           int    `json:"parentId"`           // 父类目 ID
	ExpandCatType      int    `json:"expandCatType"`      // 扩展类目类型 1=服饰，其他的非服饰
}
