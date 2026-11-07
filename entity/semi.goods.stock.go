package entity

// Stock 库存数量与类型
type Stock struct {
	Stock     int `json:"stock"`     // 库存数量
	StockType int `json:"stockType"` // 库存类型：0-自有普通库存；1-自有预售库存
}

// SemiGoodsStock 半托管商品库存
type SemiGoodsStock struct { // 来源：Temu temu.goods.stock.query
	GoodsId          int64          `json:"goodsId"`          // 商品 ID
	SkuStockInfoList []SKuStockInfo `json:"skuStockInfoList"` // SKU 库存信息列表
}

type SelfPreSaleStock struct {
	Stock
	PreSaleEndTime         int `json:"preSaleEndTime"`         // 预售结束时间
	SkuPreSaleStockLimit   int `json:"skuPreSaleStockLimit"`   // 单个 SKU 的预售库存上限
	MallPreSaleSkuNumLimit int `json:"mallPreSaleSkuNumLimit"` // 店铺预售 SKU 数量上限
}

type SKuStockInfo struct {
	OutSkuSn          string           `json:"outSkuSn"`          // 外部 SKU 编码
	SelfOrdinaryStock Stock            `json:"selfOrdinaryStock"` // 自有普通库存信息
	CertOrdinaryStock Stock            `json:"certOrdinaryStock"` // cert 普通库存信息
	SkuId             int64            `json:"skuId"`             // SKU ID
	SelfPreSaleStock  SelfPreSaleStock `json:"selfPreSaleStock"`  // 自有预售库存信息
}

// SemiGoodsStockUpdateResult 半托管商品库存更新结果
type SemiGoodsStockUpdateResult struct { // 来源：Temu bg.local.goods.stock.edit
	GoodsId                    int64                    `json:"goodsId"`                    // 商品 ID
	SkuStockEditStatusInfoList []SkuStockEditStatusInfo `json:"skuStockEditStatusInfoList"` // 各 SKU 库存修改状态
}

// SkuStockEditStatusInfo SKU 库存编辑状态信息
type SkuStockEditStatusInfo struct {
	SkuId           int64  `json:"skuId"`           // SKU ID
	StockEditStatus bool   `json:"stockEditStatus"` // 库存修改状态：true-修改成功
	ErrorCode       int    `json:"errorCode"`       // 修改失败时的错误码
	ErrorMsg        string `json:"errorMsg"`        // 修改失败时的错误信息
}
