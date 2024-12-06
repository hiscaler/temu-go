package entity

type SimpleCategory struct {
	CatId   int64  `json:"catId"`   // 类目 ID
	CatName string `json:"catName"` // 类目名称
}

type Category struct {
	CatID       int64  `json:"catId"`
	CatName     string `json:"catName"`
	CatType     int    `json:"catType"` //  1 是服饰，其他的非服饰
	CatLevel    int    `json:"catLevel"`
	ParentCatID int64  `json:"parentCatId"`
	IsLeaf      bool   `json:"isLeaf"` // true=叶子类目
	IsHidden    bool   `json:"isHidden"`
	HiddenType  int    `json:"hiddenType"`
}
