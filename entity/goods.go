package entity

// GoodsProperty 商品属性
type GoodsProperty struct {
	Vid              int         `json:"vid"`              // 基础属性值 ID
	ValueUnit        string      `json:"valueUnit"`        // 属性值单位
	Language         interface{} `json:"language"`         // 语种
	Pid              int         `json:"pid"`              // 属性 ID
	TemplatePid      int         `json:"templatePid"`      // 模板属性 ID
	NumberInputValue string      `json:"numberInputValue"` // 数值录入
	PropValue        string      `json:"propValue"`        // 基础属性值
	ValueExtendInfo  string      `json:"valueExtendInfo"`  // 属性值扩展信息
	PropName         string      `json:"propName"`         // 引用属性名
	RefPid           int         `json:"refPid"`           // 引用属性 ID
}

// GoodsSkuSummary 商品 SKU 描叙
type GoodsSkuSummary struct {
	ProductSkuId        int64  `json:"productSkuId"` // 货品 SKU Id
	ExtCode             string `json:"extCode"`
	ProductSkuWhExtAttr struct {
		ProductSkuWeight struct {
			Value int `json:"value"`
		} `json:"productSkuWeight"`
		ProductSkuWmsVolume     interface{} `json:"productSkuWmsVolume"`
		ProductSkuBarCodes      interface{} `json:"productSkuBarCodes"`
		ProductSkuSubSellMode   interface{} `json:"productSkuSubSellMode"`
		ProductSkuSensitiveAttr struct {
			SensitiveTypes []interface{} `json:"sensitiveTypes"`
			IsSensitive    int           `json:"isSensitive"`
		} `json:"productSkuSensitiveAttr"`
		ProductSkuFragileLabels    interface{} `json:"productSkuFragileLabels"`
		ProductSkuNewSensitiveAttr struct {
			Force2NormalTypes interface{} `json:"force2NormalTypes"`
			SensitiveList     []int       `json:"sensitiveList"`
			IsForce2Normal    bool        `json:"isForce2Normal"`
		} `json:"productSkuNewSensitiveAttr"`
		ProductSkuVolumeLabel interface{} `json:"productSkuVolumeLabel"`
		ProductSkuWmsWeight   interface{} `json:"productSkuWmsWeight"`
		ProductSkuVolume      struct {
			Len    int `json:"len"`    // 长
			Width  int `json:"width"`  // 宽
			Height int `json:"height"` // 高
		} `json:"productSkuVolume"`
		ProductSkuSensitiveLimit interface{} `json:"productSkuSensitiveLimit"`
		ProductSkuWmsVolumeLabel interface{} `json:"productSkuWmsVolumeLabel"`
	} `json:"productSkuWhExtAttr"`
	VirtualStock       interface{} `json:"virtualStock"`
	ProductSkuSpecList []struct {
		SpecId         int    `json:"specId"`
		ParentSpecName string `json:"parentSpecName"`
		SpecName       string `json:"specName"`
		ParentSpecId   int    `json:"parentSpecId"`
	} `json:"productSkuSpecList"`
}

type Goods struct {
	ProductProperties []GoodsProperty `json:"productProperties"` // 货品普通属性
	ProductId         int             `json:"productId"`         // 货品 ID
	ProductJitMode    struct {
		QuickSellAgtSignStatus interface{} `json:"quickSellAgtSignStatus"` // 快速售卖协议签署状态 0-未签署 1-已签署
		MatchJitMode           bool        `json:"matchJitMode"`           // 是否 JIT 模式
	} `json:"productJitMode"` // 货品 JIT 模式信息
	ProductSkuSummaries      []GoodsSkuSummary `json:"productSkuSummaries"`
	ProductName              string            `json:"productName"`
	CreatedAt                int64             `json:"createdAt"`
	ProductSemiManaged       interface{}       `json:"productSemiManaged"`
	IsSupportPersonalization bool              `json:"isSupportPersonalization"`
	ExtCode                  string            `json:"extCode"` // 货品 SKC 外部编码
	LeafCat                  struct {
		CatId   int    `json:"catId"`   // 类目 ID
		CatName string `json:"catName"` // 类目名称
	} `json:"leafCat"` // 叶子类目
	Categories struct {
		Cat8 struct {
			CatId   int    `json:"catId"`
			CatName string `json:"catName"`
		} `json:"cat8"`
		Cat9 struct {
			CatId   int    `json:"catId"`
			CatName string `json:"catName"`
		} `json:"cat9"`
		Cat6 struct {
			CatId   int    `json:"catId"`
			CatName string `json:"catName"`
		} `json:"cat6"`
		Cat7 struct {
			CatId   int    `json:"catId"`
			CatName string `json:"catName"`
		} `json:"cat7"`
		Cat4 struct {
			CatId   int    `json:"catId"`
			CatName string `json:"catName"`
		} `json:"cat4"`
		Cat5 struct {
			CatId   int    `json:"catId"`
			CatName string `json:"catName"`
		} `json:"cat5"`
		Cat2 struct {
			CatId   int    `json:"catId"`
			CatName string `json:"catName"`
		} `json:"cat2"`
		Cat3 struct {
			CatId   int    `json:"catId"`
			CatName string `json:"catName"`
		} `json:"cat3"`
		Cat10 struct {
			CatId   int    `json:"catId"`
			CatName string `json:"catName"`
		} `json:"cat10"`
		Cat1 struct {
			CatId   int    `json:"catId"`
			CatName string `json:"catName"`
		} `json:"cat1"`
		LeafCat interface{} `json:"leafCat"`
	} `json:"categories"`
	ProductSkcId    int64  `json:"productSkcId"`
	MatchSkcJitMode bool   `json:"matchSkcJitMode"`
	MainImageUrl    string `json:"mainImageUrl"`
}
