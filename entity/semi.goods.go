package entity

// SemiGoods 半托管商品
type SemiGoods struct {
	MarketPrice           int      `json:"marketPrice"`
	GoodsId               int      `json:"goodsId"`
	OutSkuSnList          []string `json:"outSkuSnList"`
	GoodsStatusChangeTime string   `json:"goodsStatusChangeTime"`
	CrtTime               int      `json:"crtTime"`
	SkuInfoList           []struct {
		SpecList []struct {
			SpecId         int    `json:"specId"`
			ParentSpecName string `json:"parentSpecName"`
			ParentSpecId   int    `json:"parentSpecId"`
			SpecName       string `json:"specName"`
		} `json:"specList"`
		Stock int    `json:"stock"`
		SkuId int    `json:"skuId"`
		SkuSn string `json:"skuSn"`
	} `json:"skuInfoList"`
	GoodsShowSubStatus   int    `json:"goodsShowSubStatus"`
	LowTrafficTag        int    `json:"lowTrafficTag"`
	OutGoodsSn           string `json:"outGoodsSn"`
	Price                string `json:"price"`
	Currency             string `json:"currency"`
	ThumbUrl             string `json:"thumbUrl"`
	GoodsName            string `json:"goodsName"`
	SkuIdList            []int  `json:"skuIdList"`
	Status4VO            int    `json:"status4VO"`
	Quantity             int    `json:"quantity"`
	SubStatus4VO         int    `json:"subStatus4VO"`
	RestrictedTrafficTag int    `json:"restrictedTrafficTag"`
	CostTemplateId       string `json:"costTemplateId"`
	ShipmentLimitSecond  int    `json:"shipmentLimitSecond"`
	CatId                int    `json:"catId"`
	SpecName             string `json:"specName"`
	TrademarkId          int    `json:"trademarkId"`
	BrandId              int    `json:"brandId"`
	RetailPrice          struct {
		Amount   string `json:"amount"`
		Currency string `json:"currency"`
	} `json:"retailPrice"`
	ListPrice struct {
		Amount   string `json:"amount"`
		Currency string `json:"currency"`
	} `json:"listPrice"`
}
