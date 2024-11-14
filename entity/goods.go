package entity

import "gopkg.in/guregu/null.v4"

// GoodsProperty 商品属性
type GoodsProperty struct {
	Vid              int         `json:"vid"`              // 基础属性值 ID
	ValueUnit        string      `json:"valueUnit"`        // 属性值单位
	Language         null.String `json:"language"`         // 语种
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
		ProductSkuWmsVolume     any `json:"productSkuWmsVolume"`
		ProductSkuBarCodes      any `json:"productSkuBarCodes"`
		ProductSkuSubSellMode   any `json:"productSkuSubSellMode"`
		ProductSkuSensitiveAttr struct {
			SensitiveTypes []any `json:"sensitiveTypes"`
			IsSensitive    int   `json:"isSensitive"`
		} `json:"productSkuSensitiveAttr"`
		ProductSkuFragileLabels    any `json:"productSkuFragileLabels"`
		ProductSkuNewSensitiveAttr struct {
			Force2NormalTypes any   `json:"force2NormalTypes"`
			SensitiveList     []int `json:"sensitiveList"`
			IsForce2Normal    bool  `json:"isForce2Normal"`
		} `json:"productSkuNewSensitiveAttr"`
		ProductSkuVolumeLabel any `json:"productSkuVolumeLabel"`
		ProductSkuWmsWeight   int `json:"productSkuWmsWeight"`
		ProductSkuVolume      struct {
			Len    int `json:"len"`    // 长
			Width  int `json:"width"`  // 宽
			Height int `json:"height"` // 高
		} `json:"productSkuVolume"`
		ProductSkuSensitiveLimit any `json:"productSkuSensitiveLimit"`
		ProductSkuWmsVolumeLabel any `json:"productSkuWmsVolumeLabel"`
	} `json:"productSkuWhExtAttr"`
	VirtualStock       int             `json:"virtualStock"`
	ProductSkuSpecList []Specification `json:"productSkuSpecList"`
}

type GoodsCategory struct {
	CatId   int    `json:"catId"`   // 类目 ID
	CatName string `json:"catName"` // 类目名称
}

type Goods struct {
	ProductProperties []GoodsProperty `json:"productProperties"` // 货品普通属性
	ProductId         int64           `json:"productId"`         // 货品 ID
	ProductJitMode    struct {
		QuickSellAgtSignStatus null.Int `json:"quickSellAgtSignStatus"` // 快速售卖协议签署状态 0-未签署 1-已签署
		MatchJitMode           bool     `json:"matchJitMode"`           // 是否 JIT 模式
	} `json:"productJitMode"` // 货品 JIT 模式信息
	ProductSkuSummaries      []GoodsSkuSummary `json:"productSkuSummaries"`
	ProductName              string            `json:"productName"`
	CreatedAt                int64             `json:"createdAt"`
	ProductSemiManaged       any               `json:"productSemiManaged"`
	IsSupportPersonalization bool              `json:"isSupportPersonalization"`
	ExtCode                  string            `json:"extCode"` // 货品 SKC 外部编码
	LeafCat                  GoodsCategory     `json:"leafCat"` // 叶子类目
	Categories               struct {
		Cat1    GoodsCategory `json:"cat1"`
		Cat2    GoodsCategory `json:"cat2"`
		Cat3    GoodsCategory `json:"cat3"`
		Cat4    GoodsCategory `json:"cat4"`
		Cat5    GoodsCategory `json:"cat5"`
		Cat6    GoodsCategory `json:"cat6"`
		Cat7    GoodsCategory `json:"cat7"`
		Cat8    GoodsCategory `json:"cat8"`
		Cat9    GoodsCategory `json:"cat9"`
		Cat10   GoodsCategory `json:"cat10"`
		LeafCat any           `json:"leafCat"`
	} `json:"categories"`
	ProductSkcId    int64  `json:"productSkcId"`
	MatchSkcJitMode bool   `json:"matchSkcJitMode"`
	MainImageUrl    string `json:"mainImageUrl"`
}
