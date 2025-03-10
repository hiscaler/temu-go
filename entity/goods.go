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
		ProductSkuWmsWeight   struct {
			WmsCollectionSourceType any `json:"wmsCollectionSourceType"`
			Value                   int `json:"value"`
		} `json:"productSkuWmsWeight"`
		ProductSkuVolume struct {
			Len    int `json:"len"`    // 长
			Width  int `json:"width"`  // 宽
			Height int `json:"height"` // 高
		} `json:"productSkuVolume"`
		ProductSkuSensitiveLimit any `json:"productSkuSensitiveLimit"`
		ProductSkuWmsVolumeLabel any `json:"productSkuWmsVolumeLabel"`
	} `json:"productSkuWhExtAttr"`
	VirtualStock          int             `json:"virtualStock"`
	ProductSkuSpecList    []Specification `json:"productSkuSpecList"`
	ProductSkuSaleExtAttr struct {
		ProductSkuShippingMode int `json:"productSkuShippingMode"` // 1：卖家自发货、2：认证仓发货
	} `json:"productSkuSaleExtAttr"` // 货品 sku 销售域扩展属性
}

// Goods 商品
type Goods struct {
	ProductProperties []GoodsProperty `json:"productProperties"` // 货品普通属性
	ProductId         int64           `json:"productId"`         // 货品 ID
	ProductJitMode    struct {
		QuickSellAgtSignStatus null.Int `json:"quickSellAgtSignStatus"` // 快速售卖协议签署状态（0：未签署、1：已签署）
		SignLatestJitVersion   bool     `json:"signLatestJitVersion"`   // 是否签署最新版本 JIT 预售协议
		MatchJitMode           bool     `json:"matchJitMode"`           // 是否 JIT 模式
	} `json:"productJitMode"` // 货品 JIT 模式信息
	ProductSkuSummaries      []GoodsSkuSummary `json:"productSkuSummaries"`
	ProductName              string            `json:"productName"`
	CreatedAt                int64             `json:"createdAt"`
	ProductSemiManaged       any               `json:"productSemiManaged"`
	IsSupportPersonalization bool              `json:"isSupportPersonalization"`
	ExtCode                  string            `json:"extCode"` // 货品 SKC 外部编码
	LeafCat                  SimpleCategory    `json:"leafCat"` // 叶子类目
	Categories               struct {
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
		LeafCat any            `json:"leafCat"`
	} `json:"categories"`
	ProductSkcId    int64  `json:"productSkcId"`
	MatchSkcJitMode bool   `json:"matchSkcJitMode"`
	MainImageUrl    string `json:"mainImageUrl"`
}

// GoodsDetail 商品详情
type GoodsDetail struct {
	ProductId   int64  `json:"productId"`   // 商品 ID
	ProductName string `json:"productName"` // 商品名称
	Categories  struct {
		CatType null.Int       `json:"catType"` // 类目类型 (0: 未分类, 1: 服饰)
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
		LeafCat SimpleCategory `json:"leafCat"` // 叶子类目
	} `json:"categories"` // 货品类目
	GoodsLayerDecorationList []struct {
		FloorId     null.Int `json:"floorId"`  // 楼层 id（null:新增，否则为更新）
		Type        string   `json:"type"`     // 组件类型 type
		Priority    int      `json:"priority"` // 楼层排序
		Lang        string   `json:"lang"`     // 语言类型
		Key         string   `json:"key"`      // 楼层类型的 key
		ContentList []struct {
			ImgUrl            string `json:"imgUrl"` // 图片地址--通用
			Width             int    `json:"width"`  // 图片宽度--通用
			Height            int    `json:"height"` // 图片高度--通用
			Text              string `json:"text"`   // 文字信息--文字模块
			TextModuleDetails struct {
				BackgroundColor string `json:"backgroundColor"` // 背景颜色
				FontFamily      int    `json:"fontFamily"`      // 字体类型
				FontSize        int    `json:"fontSize"`        // 文字模块字体大小
				FontColor       string `json:"fontColor"`       // 文字颜色
				Align           string `json:"align"`           // 文字对齐方式，left--左对齐；right--右对齐；center--居中；justify--两端对齐
			} `json:"textModuleDetails"` // 文字模块详情
		} `json:"contentList"` // 楼层内容
	} `json:"goodsLayerDecorationList"`
}
