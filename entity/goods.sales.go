package entity

import "gopkg.in/guregu/null.v4"

type GoodsSales struct {
	ProductName           string `json:"productName"` // 货品名称
	SkuQuantityDetailList []struct {
		AvailableSaleDaysFromInventory null.Float `json:"availableSaleDaysFromInventory"` // 库存可售天数，在途+在仓的库存总天数
		SkuExtCode                     string     `json:"skuExtCode"`                     // sku货号
		ClassName                      string     `json:"className"`                      // 规格名称
		LackQuantity                   int        `json:"lackQuantity"`                   // 缺货数量
		LastSevenDaysSaleVolume        int        `json:"lastSevenDaysSaleVolume"`        // 近7天销量
		ProductSkuId                   int64      `json:"productSkuId"`                   // productSkuId
		LastThirtyDaysSaleVolume       int        `json:"lastThirtyDaysSaleVolume"`       // 近 30 天销量
		AvailableSaleDays              null.Float `json:"availableSaleDays"`              // 可售天数，待发+在途+在仓的库存总天数
		TodaySaleVolume                int        `json:"todaySaleVolume"`                // 今日销量
		AdviceQuantity                 null.Int   `json:"adviceQuantity"`                 // 建议下单量
		InventoryNumInfo               struct {
			WaitOnShelfNum                   int `json:"waitOnShelfNum"`           // 待上架库存
			WaitDeliveryInventoryNum         int `json:"waitDeliveryInventoryNum"` // 待发货库存
			WarehouseInventoryNum            int `json:"warehouseInventoryNum"`    // 仓内可用库存
			WaitApproveInventoryNum          int `json:"waitApproveInventoryNum"`  // 待审核备货库存
			UnavailableWarehouseInventoryNum int `json:"unavailableWarehouseInventoryNum"`
			WaitReceiveNum                   int `json:"waitReceiveNum"` //  待收货库存
		} `json:"inventoryNumInfo"`                                               // 库存信息
		WarehouseAvailableSaleDays null.Float `json:"warehouseAvailableSaleDays"` // 仓内库存可售天数:保留一位小数
	} `json:"skuQuantityDetailList"` // sku维度数量信息
	ProductSkcId         int64       `json:"productSkcId"`
	SkcExtCode           null.String `json:"skcExtCode"` // skc货号
	ProductId            int64       `json:"productId"`
	InBlackList          bool        `json:"inBlackList"` // 是否在备货黑名单内，在的话禁止备货
	SkuQuantityTotalInfo struct {
		ProductSkuId     null.Int `json:"productSkuId"` // skuId，为null
		InventoryNumInfo struct {
			WaitDeliveryInventoryNum         int `json:"waitDeliveryInventoryNum"`         // 待发货库存
			WarehouseInventoryNum            int `json:"warehouseInventoryNum"`            // 仓库可用库存
			UnavailableWarehouseInventoryNum int `json:"unavailableWarehouseInventoryNum"` // 仓库暂不可用库存
			WaitReceiveNum                   int `json:"waitReceiveNum"`                   // 待收货库存
		} `json:"inventoryNumInfo"` // 库存汇总信息，SKU 维度
	} `json:"skuQuantityTotalInfo"` // sku汇总库存数据
	SupplyStatusRemark     null.String `json:"supplyStatusRemark"`
	OnSalesDurationOffline int         `json:"onSalesDurationOffline"` // 加入站点天数
	ProductSkcPicture      string      `json:"productSkcPicture"`      // 货品图片
	Category               string      `json:"category"`
	Mark                   float64     `json:"mark"`
}
