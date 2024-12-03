package entity

// GoodsTopSellingSoldOut 爆款售罄商品
type GoodsTopSellingSoldOut struct {
	SellOutProductId int64  `json:"sellOutProductId"` // 售罄货品 id
	ProductName      string `json:"productName"`      // 售罄货品名称
	ProductPicture   string `json:"productPicture"`   // 售罄货品主图
	Categories       struct {
		CatType int            `json:"catType"` // 类目类型 (0: 未分类, 1: 服饰)
		Cat1    SimpleCategory `json:"cat1"`
		Cat2    SimpleCategory `json:"cat2"`
		Cat3    SimpleCategory `json:"cat3"`
		Cat4    SimpleCategory `json:"cat4"`
		Cat5    SimpleCategory `json:"cat5"`
		Cat6    SimpleCategory `json:"cat6"`
		Cat7    SimpleCategory `json:"cat7"`
		Cat8    SimpleCategory `json:"cat8"`
		Cat9    SimpleCategory `json:"cat9"`
		Cat10   SimpleCategory `json:"cat10"`
	} `json:"categories"` // 售罄货品类目
}
