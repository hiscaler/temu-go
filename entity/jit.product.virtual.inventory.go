package entity

import "gopkg.in/guregu/null.v4"

// JitProductVirtualInventory JIT 货品虚拟库存
type JitProductVirtualInventory struct {
	ProductSkuId     int64    `json:"productSkuId"`     // 货品 SKU ID
	SkuStockQuantity null.Int `json:"skuStockQuantity"` // 货品 SKU 虚拟库存, 不允许查看时返回 null
}
