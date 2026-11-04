package entity

// SemiGoods 半托管商品
// 来源：Temu bg.local.goods.list.query
type SemiGoods struct {
	MarketPrice           int64    `json:"marketPrice"`           // 市场价或建议零售价
	GoodsId               int64    `json:"goodsId"`               // 商品唯一标识
	OutSkuSnList          []string `json:"outSkuSnList"`          // 商品关联的外部 SKU 编码列表
	GoodsStatusChangeTime string   `json:"goodsStatusChangeTime"` // 商品状态变更时间，时间戳格式
	CrtTime               int64    `json:"crtTime"`               // 商品创建时间，Unix 时间戳（秒）
	SkuInfoList           []struct {
		SpecList []struct {
			SpecId         int64  `json:"specId"`         // 规格值 ID
			ParentSpecName string `json:"parentSpecName"` // 父规格名称
			ParentSpecId   int64  `json:"parentSpecId"`   // 父规格 ID
			SpecName       string `json:"specName"`       // 规格值名称
		} `json:"specList"` // SKU 规格值列表
		Stock int    `json:"stock"` // SKU 库存数量
		SkuId int64  `json:"skuId"` // SKU 唯一标识
		SkuSn string `json:"skuSn"` // 外部 SKU 编码
	} `json:"skuInfoList"` // SKU 信息列表
	GoodsShowSubStatus   int            `json:"goodsShowSubStatus"`   // 商品子状态（新版字段）
	LowTrafficTag        int            `json:"lowTrafficTag"`        // 低流量标签：1-低流量，2-非低流量
	OutGoodsSn           string         `json:"outGoodsSn"`           // 外部商品编码
	Price                string         `json:"price"`                // 商品售价
	Currency             string         `json:"currency"`             // 商品售价币种
	ThumbUrl             string         `json:"thumbUrl"`             // 商品缩略图地址
	GoodsName            string         `json:"goodsName"`            // 商品标题
	SkuIdList            []int64        `json:"skuIdList"`            // 商品关联的 SKU ID 列表
	Status4VO            int            `json:"status4VO"`            // 商品状态：1-在售，4-未发布，5-草稿，6-已删除
	Quantity             int            `json:"quantity"`             // 商品库存数量
	SubStatus4VO         int            `json:"subStatus4VO"`         // 商品子状态
	RestrictedTrafficTag int            `json:"restrictedTrafficTag"` // 流量限制标签：1-限制流量，2-未限制流量
	CostTemplateId       string         `json:"costTemplateId"`       // 商品可用配送选项 ID，多个 ID 以英文逗号分隔
	ShipmentLimitSecond  int64          `json:"shipmentLimitSecond"`  // 从接收订单到可发货的时长，单位为秒
	CatId                int64          `json:"catId"`                // 商品类目 ID
	SpecName             string         `json:"specName"`             // 商品规格名称
	TrademarkId          int64          `json:"trademarkId"`          // 商品商标 ID
	BrandId              int64          `json:"brandId"`              // 商品品牌 ID
	RetailPrice          SemiGoodsPrice `json:"retailPrice"`          // 商品售价及币种
	ListPrice            SemiGoodsPrice `json:"listPrice"`            // 商品市场价或建议零售价及币种
}

// SemiGoodsPrice 半托管商品价格
type SemiGoodsPrice struct {
	Amount   string `json:"amount"`   // 价格金额
	Currency string `json:"currency"` // 价格币种
}
