package entity

import "gopkg.in/guregu/null.v4"

type GoodsQuantity struct {
	ProductSkuId     int64    `json:"productSkuId"`     // 货品SKUId
	SkuStockQuantity null.Int `json:"skuStockQuantity"` // 货品SKU虚拟库存, 不允许查看时返回null
	WarehouseId      string   `json:"warehouseId"`      // 仓库ID，货品SKUId维度数据，欧洲地区支持分仓库存
	ShippingMode     int      `json:"shippingMode"`     // 发货模式：1-卖家自发货，2-合作对接仓托管
	TempLockQuantity null.Int `json:"tempLockQuantity"` // 未支付的库存数量（这部分库存可能存在消费者取消订单后的库存返增）
}
