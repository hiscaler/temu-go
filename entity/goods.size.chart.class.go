package entity

// GoodsSizeChartClass 尺码分类对象
type GoodsSizeChartClass struct {
	CatId           int `json:"catId"`           // 叶子类目ID
	ClassId         int `json:"classId"`         // 分类ID
	ClassType       int `json:"classType"`       // 分类类型，0：单尺码表，1-多尺码表
	ParentClassId   int `json:"parentClassId"`   // 父分类ID
	RelatedClassIds int `json:"relatedClassIds"` // 关联分类ID
}
