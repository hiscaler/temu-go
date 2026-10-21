package entity

type Stock struct {
	Stock     int `json:"stock"`
	StockType int `json:"stockType"`
}

// SemiGoodsStock 半托管商品库存
type SemiGoodsStock struct {
	GoodsId          int `json:"goodsId"`
	SkuStockInfoList []struct {
		OutSkuSn          string `json:"outSkuSn"`
		SelfOrdinaryStock Stock  `json:"selfOrdinaryStock"`
		CertOrdinaryStock Stock  `json:"certOrdinaryStock"`
		SkuId             int    `json:"skuId"`
		SelfPreSaleStock  struct {
			Stock
			PreSaleEndTime         int `json:"preSaleEndTime"`
			SkuPreSaleStockLimit   int `json:"skuPreSaleStockLimit"`
			MallPreSaleSkuNumLimit int `json:"mallPreSaleSkuNumLimit"`
		} `json:"selfPreSaleStock"`
	} `json:"skuStockInfoList"`
}

// SemiGoodsStockUpdateResult 半托管商品库存更新结果
type SemiGoodsStockUpdateResult struct {
	GoodsId                    int64 `json:"goodsId"`
	SkuStockEditStatusInfoList []struct {
		SkuId           int64  `json:"skuId"`
		StockEditStatus bool   `json:"stockEditStatus"` // Stock Edit Status, if true means successful stock modification
		ErrorCode       int    `json:"errorCode"`
		ErrorMsg        string `json:"errorMsg"`
	} `json:"skuStockEditStatusInfoList"`
}
