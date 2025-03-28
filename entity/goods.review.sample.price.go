package entity

type GoodsReviewSamplePrice struct {
	ProductSkuIdList     []int64  `json:"productSkuIdList"`     // 货品 skuId
	PriceCurrency        string   `json:"priceCurrency"`        // 申报价格币种
	SupplyPrice          int      `json:"supplyPrice"`          // 申报价格（分）
	OrderId              int64    `json:"orderId"`              // 核价单id
	SuggestPriceCurrency string   `json:"suggestPriceCurrency"` // 建议价格币种
	SuggestSupplyPrice   int      `json:"suggestSupplyPrice"`   // 建议价格（分）
	OrderStatus          int      `json:"orderStatus"`          // 核价单的状态. 可选值含义说明:[0:待核价;1:待供应商确认;2:核价通过;3:核价驳回;4:废弃;5:价格同步中;]
	SiteIds              []int    `json:"siteIds"`              // 站点 id
	CanBargain           bool     `json:"canBargain"`           // 是否可重新报价
	SiteNameList         []string `json:"siteNameList"`         // 站点名称
}
