package temu

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
	"gopkg.in/guregu/null.v4"
)

type semiGoodsService struct {
	service
	Image semiGoodsImageService
}

type SemiGoodsQueryParams struct {
	normal.ParameterWithPager
	OrderField                string  `json:"orderField,omitempty"`
	OrderType                 int     `json:"orderType,omitempty"`
	GoodsSearchType           int     `json:"goodsSearchType,omitempty"`
	SearchText                string  `json:"searchText,omitempty"`
	StatusFilterType          int     `json:"statusFilterType,omitempty"`
	CrtFrom                   int64   `json:"crtFrom,omitempty"`
	CrtTo                     int64   `json:"crtTo,omitempty"`
	GoodsIdList               []int64 `json:"goodsIdList,omitempty"`
	CatIdList                 []int64 `json:"catIdList,omitempty"`
	GoodsStatusFilterType     int     `json:"goodsStatusFilterType"`
	GoodsSubStatusFilterType  int     `json:"goodsSubStatusFilterType"`
	GoodsStatusChangeTimeFrom int64   `json:"goodsStatusChangeTimeFrom"`
	GoodsStatusChangeTimeTo   int64   `json:"goodsStatusChangeTimeTo"`
	GoodsSearchTags           []int   `json:"goodsSearchTags"`
}

// validate 校验半托管商品查询参数。
func (m SemiGoodsQueryParams) validate() error {
	return nil
}

// Query 商品查询
// https://partner-us.temu.com/documentation?menu_code=fb16b05f7a904765aac4af3a24b87d4a&sub_menu_code=d2a836cf1711473ba1f83597a1b52fb0
func (s semiGoodsService) Query(ctx context.Context, params SemiGoodsQueryParams) (items []entity.SemiGoods, total, totalPages int, isLastPage bool, err error) {
	params.Page = params.TidyPager().Page
	params.OmitPage()
	if err = params.validate(); err != nil {
		return items, 0, 0, true, invalidInput(err)
	}

	var result = struct {
		normal.Response
		Result struct {
			PageNo    int                `json:"pageNo"`
			Total     int                `json:"total"`
			GoodsList []entity.SemiGoods `json:"goodsList"`
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(params).
		SetResult(&result).
		Post("bg.local.goods.list.query")
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	items = result.Result.GoodsList
	total, totalPages, isLastPage = parseResponseTotal(params.Page, params.PageSize, result.Result.Total)
	return
}

// SemiGoodsCreateRequest temu.local.goods.v3.add 添加商品请求参数
type SemiGoodsCreateRequest struct {
	Language   string                     `json:"language"`   // 语言，可选
	GoodsBasic SemiGoodsCreateGoodsBasic  `json:"goodsBasic"` // 商品基础信息，必填
	Attributes []SemiGoodsCreateAttribute `json:"attributes"` // 商品属性，可填写基础属性、治理属性、合规属性、品牌及二手商品信息，可选
	SkuList    []SemiGoodsCreateSku       `json:"skuList"`    // SKU 列表，必填，至少包含一个 SKU
}

// SemiGoodsCreateGoodsBasic temu.local.goods.v3.add 商品基础信息
type SemiGoodsCreateGoodsBasic struct {
	ExternalGoodsId      string   `json:"externalGoodsId"`                // 外部平台商品编码，用于关联 Temu 商品，同一店铺内必须唯一，必填
	GoodsName            string   `json:"goodsName"`                      // 商品名称，必填
	ExtCatName           string   `json:"extCatName,omitempty"`           // 外部类目名称，可提高类目匹配准确度
	AssignedPrimaryCatId int64    `json:"assignedPrimaryCatId,omitempty"` // 指定的 Temu 一级类目 ID
	GoodsDesc            string   `json:"goodsDesc,omitempty"`            // 商品描述
	GoodsCarouselImage   []string `json:"goodsCarouselImage,omitempty"`   // 商品轮播图
	DetailImage          []string `json:"detailImage,omitempty"`          // 商品详情图
	ProductType          int      `json:"productType,omitempty"`          // 商品类型：1-普通，2-定制，3-MTO，4-二手
	ShipmentLimitDay     int      `json:"shipmentLimitDay,omitempty"`     // 承诺发货生产时间，单位：天
	BulletPoints         []string `json:"bulletPoints,omitempty"`         // 商品卖点
	CostTemplate         string   `json:"costTemplate,omitempty"`         // 运费模板 ID 或名称
}

// SemiGoodsCreateAttribute temu.local.goods.v3.add 商品属性
type SemiGoodsCreateAttribute struct {
	Name  string   `json:"name"`  // 属性名称
	Value []string `json:"value"` // 属性值
}

// SemiGoodsCreateSku temu.local.goods.v3.add SKU 信息
type SemiGoodsCreateSku struct {
	ExternalSkuId string                        `json:"externalSkuId"`        // 外部 SKU 编码，同一店铺内必须唯一，必填
	Images        []string                      `json:"images"`               // SKU 轮播图，必填
	Price         SemiGoodsCreateSkuPrice       `json:"price"`                // 价格信息，必填
	Quantity      int64                         `json:"quantity"`             // 常规库存数量，必填
	PackageInfo   SemiGoodsCreateSkuPackageInfo `json:"packageInfo"`          // 商品包装信息，必填
	Variations    []SemiGoodsCreateSkuVariation `json:"variations"`           // SKU 规格，至少包含一个，必填
	BarCode       *SemiGoodsCreateSkuBarCode    `json:"barCode,omitempty"`    // 外部商品编码，可选
	References    *SemiGoodsCreateSkuReferences `json:"references,omitempty"` // 参考信息，可选
}

// SemiGoodsCreateSkuPrice temu.local.goods.v3.add SKU 价格信息
type SemiGoodsCreateSkuPrice struct {
	BasePrice SemiGoodsCreateSkuPriceValue  `json:"basePrice"`           // 基础价格，必填
	ListPrice *SemiGoodsCreateSkuPriceValue `json:"listPrice,omitempty"` // 建议零售价，必须高于基础价格，可选
}

// SemiGoodsCreateSkuPriceValue temu.local.goods.v3.add 价格值
type SemiGoodsCreateSkuPriceValue struct {
	Amount   string `json:"amount"`   // 金额，必填
	Currency string `json:"currency"` // 货币，必填
}

// SemiGoodsCreateSkuPackageInfo temu.local.goods.v3.add SKU 包装信息
type SemiGoodsCreateSkuPackageInfo struct {
	Weight string `json:"weight"` // 重量，美国使用 lb，其他国家使用 g，必填
	Length string `json:"length"` // 长度，美国使用 in，其他国家使用 cm，必填
	Width  string `json:"width"`  // 宽度，美国使用 in，其他国家使用 cm，必填
	Height string `json:"height"` // 高度，美国使用 in，其他国家使用 cm，必填
}

// SemiGoodsCreateSkuVariation temu.local.goods.v3.add SKU 规格
type SemiGoodsCreateSkuVariation struct {
	Name  string `json:"name"`  // 父规格名称，例如 color、size，必填
	Value string `json:"value"` // 规格值，例如 red、xl，必填
}

// SemiGoodsCreateSkuBarCode temu.local.goods.v3.add SKU 外部商品编码
type SemiGoodsCreateSkuBarCode struct {
	BarCodeType string   `json:"barCodeType"` // 编码类型：EAN、UPC、ISBN、GTIN-14
	BarCodeId   []string `json:"barCodeId"`   // 编码值，ISBN 支持传入多个
}

// SemiGoodsCreateSkuReferences temu.local.goods.v3.add SKU 参考信息
type SemiGoodsCreateSkuReferences struct {
	Code *SemiGoodsCreateSkuReferenceCode `json:"code,omitempty"` // 参考编码
	Note []string                         `json:"note,omitempty"` // 其他参考信息
}

// SemiGoodsCreateSkuReferenceCode temu.local.goods.v3.add SKU 参考编码
type SemiGoodsCreateSkuReferenceCode struct {
	Name string `json:"name"` // 编码名称
	Id   string `json:"id"`   // 编码 ID
}

// validate 校验 temu.local.goods.v3.add 添加商品请求参数。
func (m SemiGoodsCreateRequest) validate() error {
	if err := validation.ValidateStruct(&m,
		validation.Field(&m.SkuList,
			validation.Required.Error("SKU 列表不能为空"),
			validation.Min(1).Error("SKU 列表至少包含一个 SKU"),
		),
	); err != nil {
		return err
	}

	if err := m.GoodsBasic.validate(); err != nil {
		return err
	}
	for _, attribute := range m.Attributes {
		if err := attribute.validate(); err != nil {
			return err
		}
	}
	for _, sku := range m.SkuList {
		if err := sku.validate(); err != nil {
			return err
		}
	}

	return nil
}

// validate 校验商品基础信息。
func (m SemiGoodsCreateGoodsBasic) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.ExternalGoodsId, validation.Required.Error("外部商品 ID 不能为空")),
		validation.Field(&m.GoodsName, validation.Required.Error("商品名称不能为空")),
		validation.Field(&m.AssignedPrimaryCatId,
			validation.When(m.AssignedPrimaryCatId != 0, validation.Min(1).Error("指定类目 ID 必须大于 0")),
		),
		validation.Field(&m.ProductType,
			validation.When(m.ProductType != 0, validation.In(1, 2, 3, 4).Error("无效的商品类型")),
		),
		validation.Field(&m.ShipmentLimitDay,
			validation.When(m.ShipmentLimitDay != 0, validation.Min(1).Error("承诺发货生产时间必须大于 0")),
		),
		validation.Field(&m.GoodsCarouselImage,
			validation.Each(validation.Required.Error("商品轮播图不能为空")),
		),
		validation.Field(&m.DetailImage,
			validation.Each(validation.Required.Error("商品详情图不能为空")),
		),
		validation.Field(&m.BulletPoints,
			validation.Each(validation.Required.Error("商品卖点不能为空")),
		),
	)
}

// validate 校验商品属性。
func (m SemiGoodsCreateAttribute) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Required.Error("属性名称不能为空")),
		validation.Field(&m.Value,
			validation.Required.Error("属性值不能为空"),
			validation.Min(1).Error("属性值至少包含一个"),
			validation.Each(validation.Required.Error("属性值不能为空")),
		),
	)
}

// validate 校验 SKU 信息及其嵌套数据。
func (m SemiGoodsCreateSku) validate() error {
	if err := validation.ValidateStruct(&m,
		validation.Field(&m.ExternalSkuId, validation.Required.Error("外部 SKU ID 不能为空")),
		validation.Field(&m.Images,
			validation.Required.Error("SKU 图片不能为空"),
			validation.Min(1).Error("SKU 图片至少包含一张"),
			validation.Each(validation.Required.Error("SKU 图片不能为空")),
		),
		validation.Field(&m.Quantity, validation.Min(0).Error("库存数量不能为负数")),
		validation.Field(&m.Variations,
			validation.Required.Error("SKU 规格不能为空"),
			validation.Min(1).Error("SKU 规格至少包含一个"),
		),
	); err != nil {
		return err
	}

	if err := m.Price.validate(); err != nil {
		return err
	}
	if err := m.PackageInfo.validate(); err != nil {
		return err
	}
	for _, variation := range m.Variations {
		if err := variation.validate(); err != nil {
			return err
		}
	}
	if m.BarCode != nil {
		if err := m.BarCode.validate(); err != nil {
			return err
		}
	}
	if m.References != nil {
		if err := m.References.validate(); err != nil {
			return err
		}
	}

	return nil
}

// validate 校验 SKU 价格信息。
func (m SemiGoodsCreateSkuPrice) validate() error {
	if err := m.BasePrice.validate(); err != nil {
		return err
	}
	if m.ListPrice != nil {
		return m.ListPrice.validate()
	}
	return nil
}

// validate 校验价格金额和币种。
func (m SemiGoodsCreateSkuPriceValue) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Amount, validation.Required.Error("价格金额不能为空")),
		validation.Field(&m.Currency, validation.Required.Error("价格币种不能为空")),
	)
}

// validate 校验 SKU 包装信息。
func (m SemiGoodsCreateSkuPackageInfo) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Weight, validation.Required.Error("包装重量不能为空")),
		validation.Field(&m.Length, validation.Required.Error("包装长度不能为空")),
		validation.Field(&m.Width, validation.Required.Error("包装宽度不能为空")),
		validation.Field(&m.Height, validation.Required.Error("包装高度不能为空")),
	)
}

// validate 校验 SKU 规格信息。
func (m SemiGoodsCreateSkuVariation) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Required.Error("规格名称不能为空")),
		validation.Field(&m.Value, validation.Required.Error("规格值不能为空")),
	)
}

// validate 校验 SKU 条码信息。
func (m SemiGoodsCreateSkuBarCode) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.BarCodeType,
			validation.When(m.BarCodeType != "", validation.In("EAN", "UPC", "ISBN", "GTIN-14").Error("无效的条码类型")),
		),
		validation.Field(&m.BarCodeId,
			validation.Each(validation.Required.Error("条码值不能为空")),
		),
	)
}

// validate 校验 SKU 参考信息。
func (m SemiGoodsCreateSkuReferences) validate() error {
	if m.Code == nil {
		return nil
	}

	return validation.ValidateStruct(m.Code,
		validation.Field(&m.Code.Name, validation.Required.Error("参考编码名称不能为空")),
		validation.Field(&m.Code.Id, validation.Required.Error("参考编码 ID 不能为空")),
	)
}

// Create 添加商品
// https://partner-us.temu.com/documentation?menu_code=fb16b05f7a904765aac4af3a24b87d4a&sub_menu_code=f6f78d3f42ca41e59420cc4b7e3be479
func (s semiGoodsService) Create(ctx context.Context, req SemiGoodsCreateRequest) (int64, string, error) {
	if err := req.validate(); err != nil {
		return 0, "", invalidInput(err)
	}

	var result = struct {
		normal.Response
		Result struct {
			GoodsId         int64  `json:"goodsId"`
			ExternalGoodsId string `json:"externalGoodsId"`
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(req).
		SetResult(&result).
		Post("temu.local.goods.v3.add")
	if err = recheckError(resp, result.Response, err); err != nil {
		return 0, "", err
	}

	return result.Result.GoodsId, result.Result.ExternalGoodsId, nil
}

type SemiGoodsInferRequest struct {
	GoodsId        int64   `json:"goodsId"`         // 商品 ID
	Cost           float64 `json:"cost"`            // 成本
	Freight        float64 `json:"freight"`         // 运费
	ExpectedProfit float64 `json:"expected_profit"` // 期望利润
}

type SemiGoodsInferRecommend struct {
	ActivityId null.Int   `json:"activity_id"` // 活动 ID
	Roas       null.Float `json:"roas"`        // ROAS
	Profit     float64    `json:"profit"`      // 利润（申报价|活动价 - 成本 - 运费 - (申报价|活动价)/Roas）
}

type SemiGoodsInferResult struct {
	GoodsId    int64                  `json:"goods_id"`
	Recommends []SemiGoodsInferResult `json:"recommends"` // 推荐的
}

func (s semiGoodsService) Infer(ctx context.Context, req SemiGoodsInferRequest) (*SemiGoodsInferResult, error) {
	return nil, nil
}
