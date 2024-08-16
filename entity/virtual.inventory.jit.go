package entity

import "gopkg.in/guregu/null.v4"

// VirtualInventoryJit 虚拟库存 Jit
type VirtualInventoryJit struct {
	ProductSkuId     int      `json:"productSkuId"`     // 货品SKUId
	SkuStockQuantity null.Int `json:"skuStockQuantity"` // 货品SKU虚拟库存, 不允许查看时返回null
}
