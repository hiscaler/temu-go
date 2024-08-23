package entity

// GoodsLabel 商品条码
type GoodsLabel struct {
	ProductSkuSpecI18nMap Specification `json:"productSkuSpecI18nMap"` // sku规格多语言信息
	ProductSkuDTO         struct {
		ProductSkuId int    `json:"productSkuId"` // 货品 skuId
		ExtCode      string `json:"extCode"`      // sku 货号
		ProductId    int    `json:"productId"`    // 货品 Id
	} `json:"productSkuDTO"` // sku信息
	ProductLabelCodeDTO struct {
		ProductSkuId               int    `json:"productSkuId"`
		CreateTime                 int    `json:"createTime"`
		PurchaseOrderSn            string `json:"purchaseOrderSn"`
		SubPurchaseOrderSn         string `json:"subPurchaseOrderSn"`
		ProductSkcId               int    `json:"productSkcId"`
		ProductSkuPurchaseQuantity int    `json:"productSkuPurchaseQuantity"` // sku 下单件数 (仅旧版分页查询接口返回)
		LabelCode                  int    `json:"商品条码"`                       // sku 下单件数 (仅旧版分页查询接口返回)
	} `json:"ProductLabelCodeDTO"`
	ProductSkcImageList struct {
		ImageUrl  string `json:"imageUrl"`  // 图片 URL
		Language  string `json:"language"`  // 语言
		ImageType int    `json:"imageType"` // 图片类型
	} `json:"productSkcImageList"` // skc图片信息
	ProductOrigin struct {
		CountryShortName string `json:"countryShortName"` // 国家简称 (二字简码)
		CountryName      string `json:"countryName"`      // 国家名称 (英文)
	} `json:"productOrigin"` // 货品产地信息
	ProductSkcDTO struct {
		SpecIdList        []int          `json:"specIdList"`        // 主销售属性id列表
		ExtCode           string         `json:"extCode"`           // skc货号
		ProductId         int            `json:"productId"`         // 货品 Id
		ProductSkcSpec    any            `json:"productSkcSpec"`    // 主销售属性详情
		ProductSkcSpecMap map[string]any `json:"productSkcSpecMap"` // skc主销售规格Map
		ProductSkcId      string         `json:"productSkcId"`      // 货品 skcId
	} `json:"productSkcDTO"` // skc 信息
	ProductDTO struct {
		SupplierName    string `json:"supplierName"`    // 供应商名称
		LeafCatLabel    any    `json:"leafCatLabel"`    // 叶子类目标记 (使用前请与接口提供者确认是否会返回该字段)
		ProductId       int    `json:"productId"`       // 货品 ID
		ProductI18nList any    `json:"productI18nList"` // 货品多语言信息
		SourceType      int    `json:"sourceType"`      // 来源
		Categories      any    `json:"categories"`      // 类目
		ProductName     string `json:"productName"`     // 货品名称
		ProductType     int    `json:"productType"`     // 货品类型
	} `json:"productDTO"` // spu信息
	ProductSkuLabelCodeDTO struct {
		ProductSkuId int `json:"productSkuId"` // 货品 sku id
		ProductId    int `json:"productId"`    // 货品 id
		ProductSkcId int `json:"productSkcId"` // 货品 skc id
		LabelCode    int `json:"labelCode"`    // 标签条码
	} `json:"productSkuLabelCodeDTO"` // 新版货品标签条码基础信息
	ProductSkcSpecI18NMap map[string][]Specification `json:"productSkcSpecI18nMap"`
}

// CustomGoodsLabel 定制商品条码
type CustomGoodsLabel struct {
	ProductSkuSpecI18NMap map[string][]Specification `json:"productSkuSpecI18nMap"`
	ProductSkuDTO         struct {
		NumberOfPieces any `json:"numberOfPieces"`
		ProductSkuID   int `json:"productSkuId"`
		ProductID      int `json:"productId"`
		ProductSkuSpec struct {
			ProductSkuID int             `json:"productSkuId"`
			ProductID    int             `json:"productId"`
			SpecList     []Specification `json:"specList"`
			ProductSkcID int             `json:"productSkcId"`
		} `json:"productSkuSpec"`
		PieceUnitCode     any    `json:"pieceUnitCode"`
		ExtCode           string `json:"extCode"`
		ProductSkcID      int    `json:"productSkcId"`
		ThumbURL          string `json:"thumbUrl"`
		SkuClassification any    `json:"skuClassification"`
	} `json:"productSkuDTO"`
	ProductSkcImageList []struct {
		ImageURL  string `json:"imageUrl"`
		Language  string `json:"language"`
		ImageType int    `json:"imageType"`
	} `json:"productSkcImageList"`
	ProductSkcDTO struct {
		SpecIDList     []int  `json:"specIdList"`
		ExtCode        string `json:"extCode"`
		ProductSkcSpec struct {
			ProductID    int             `json:"productId"`
			SpecList     []Specification `json:"specList"`
			ProductSkcID int             `json:"productSkcId"`
		} `json:"productSkcSpec"`
		ProductID    int `json:"productId"`
		ProductSkcID int `json:"productSkcId"`
	} `json:"productSkcDTO"`
	ProductOrigin struct {
		Region1ShortName string `json:"region1ShortName"`
		Region1Name      string `json:"region1Name"`
	} `json:"productOrigin"`
	ProductSkuLabelCodeDTO struct {
		ProductSkuID int `json:"productSkuId"`
		ProductID    int `json:"productId"`
		CreateTimeTs int `json:"createTimeTs"`
		ProductSkcID int `json:"productSkcId"`
		LabelCode    int `json:"labelCode"`
	} `json:"productSkuLabelCodeDTO"`
	ProductSkcSpecI18NMap map[string][]Specification `json:"productSkcSpecI18nMap"`
}
