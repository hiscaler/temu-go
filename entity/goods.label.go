package entity

// GoodsLabel 商品条码
type GoodsLabel struct {
	ProductSkuSpecI18nMap struct {
		SpecId         int    `json:"specId"`         // 规格id
		ParentSpecName string `json:"parentSpecName"` // 父规格名称
		ParentSpecId   int    `json:"parentSpecId"`   // 父规格id
		SpecName       string `json:"specName"`       // 规格名称
	} `json:"productSkuSpecI18nMap"` // sku规格多语言信息
	ProductSkuDTO struct {
		ProductSkuId int    `json:"productSkuId"` // 货品skuId
		ExtCode      string `json:"extCode"`      // sku货号
		ProductId    int    `json:"productId"`    // 货品id
	} `json:"productSkuDTO"` // sku信息
	ProductLabelCodeDTO struct {
		ProductSkuId               int    `json:"productSkuId"`
		CreateTime                 int    `json:"createTime"`
		PurchaseOrderSn            string `json:"purchaseOrderSn"`
		SubPurchaseOrderSn         string `json:"subPurchaseOrderSn"`
		ProductSkcId               int    `json:"productSkcId"`
		ProductSkuPurchaseQuantity int    `json:"productSkuPurchaseQuantity"` // sku下单件数 (仅旧版分页查询接口返回)
		LabelCode                  int    `json:"商品条码"`                   // sku下单件数 (仅旧版分页查询接口返回)
	} `json:"ProductLabelCodeDTO"`
	ProductSkcImageList struct {
		ImageUrl  string `json:"imageUrl"`  // 图片URL
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
		ProductId         int            `json:"productId"`         // 货品Id
		ProductSkcSpec    any            `json:"productSkcSpec"`    // 主销售属性详情
		ProductSkcSpecMap map[string]any `json:"productSkcSpecMap"` // skc主销售规格Map
		ProductSkcId      string         `json:"productSkcId"`      // 货品skcId
	} `json:"productSkcDTO"` // skc 信息
	ProductDTO struct {
		SupplierName    string `json:"supplierName"`    // 供应商名称
		LeafCatLabel    any    `json:"leafCatLabel"`    // 叶子类目标记 (使用前请与接口提供者确认是否会返回该字段)
		ProductId       int    `json:"productId"`       // 货品ID
		ProductI18nList any    `json:"productI18nList"` // 货品多语言信息
		SourceType      int    `json:"sourceType"`      // 来源
		Categories      any    `json:"categories"`      // 类目
		ProductName     string `json:"productName"`     // 货品名称
		ProductType     int    `json:"productType"`     // 货品类型
	} `json:"productDTO"` // spu信息
	ProductSkuLabelCodeDTO struct {
		ProductSkuId int `json:"productSkuId"` // 货品sku id
		ProductId    int `json:"productId"`    // 货品id
		ProductSkcId int `json:"productSkcId"` // 货品skc id
		LabelCode    int `json:"labelCode"`    // 标签条码
	} `json:"productSkuLabelCodeDTO"` // 新版货品标签条码基础信息
	ProductSkcSpecI18nMap struct {
		SpecId         int    `json:"specId"`         // 规格id
		ParentSpecName string `json:"parentSpecName"` // 父规格名称
		ParentSpecId   int    `json:"parentSpecId"`   // 父规格id
		SpecName       string `json:"specName"`       // 规格名称
	} `json:"productSkcSpecI18nMap"` // skc规格多语言信息
}

// CustomGoodsLabel 定制商品条码
type CustomGoodsLabel struct {
	GoodsLabel
}
