package entity

type SimpleCategory struct {
	CatId   int64  `json:"catId"`   // 类目 ID
	CatName string `json:"catName"` // 类目名称
}

type Category struct {
	CatId       int64  `json:"catId"`   // 分类 ID
	CatName     string `json:"catName"` // 分类名称
	CatType     int    `json:"catType"` //  1 是服饰，其他的非服饰
	CatLevel    int    `json:"catLevel"`
	ParentCatId int64  `json:"parentCatId"`
	IsLeaf      bool   `json:"isLeaf"` // true=叶子类目
	IsHidden    bool   `json:"isHidden"`
	HiddenType  int    `json:"hiddenType"`
}
