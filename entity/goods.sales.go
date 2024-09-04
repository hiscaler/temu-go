package entity

import "gopkg.in/guregu/null.v4"

type GoodsSales struct {
	ProductName           string `json:"productName"`
	SkuQuantityDetailList []struct {
		AvailableSaleDaysFromInventory any    `json:"availableSaleDaysFromInventory"`
		SkuExtCode                     string `json:"skuExtCode"`
		ClassName                      string `json:"className"`
		LackQuantity                   int    `json:"lackQuantity"`
		LastSevenDaysSaleVolume        int    `json:"lastSevenDaysSaleVolume"`
		ProductSkuID                   int    `json:"productSkuId"`
		LastThirtyDaysSaleVolume       int    `json:"lastThirtyDaysSaleVolume"`
		AvailableSaleDays              any    `json:"availableSaleDays"`
		TodaySaleVolume                int    `json:"todaySaleVolume"`
		AdviceQuantity                 any    `json:"adviceQuantity"`
		InventoryNumInfo               struct {
			WaitOnShelfNum                   int `json:"waitOnShelfNum"`
			WaitDeliveryInventoryNum         int `json:"waitDeliveryInventoryNum"`
			WarehouseInventoryNum            int `json:"warehouseInventoryNum"`
			WaitApproveInventoryNum          int `json:"waitApproveInventoryNum"`
			UnavailableWarehouseInventoryNum int `json:"unavailableWarehouseInventoryNum"`
			WaitReceiveNum                   int `json:"waitReceiveNum"`
		} `json:"inventoryNumInfo"`
		WarehouseAvailableSaleDays any `json:"warehouseAvailableSaleDays"`
	} `json:"skuQuantityDetailList"`
	ProductSkcID         int    `json:"productSkcId"`
	SkcExtCode           string `json:"skcExtCode"`
	ProductID            int    `json:"productId"`
	InBlackList          bool   `json:"inBlackList"`
	SkuQuantityTotalInfo struct {
		ProductSkuID     null.Int `json:"productSkuId"`
		InventoryNumInfo struct {
			WaitDeliveryInventoryNum         int `json:"waitDeliveryInventoryNum"`
			WarehouseInventoryNum            int `json:"warehouseInventoryNum"`
			UnavailableWarehouseInventoryNum int `json:"unavailableWarehouseInventoryNum"`
			WaitReceiveNum                   int `json:"waitReceiveNum"`
		} `json:"inventoryNumInfo"`
	} `json:"skuQuantityTotalInfo"`
	SupplyStatusRemark     string  `json:"supplyStatusRemark"`
	OnSalesDurationOffline int     `json:"onSalesDurationOffline"`
	ProductSkcPicture      string  `json:"productSkcPicture"`
	Category               string  `json:"category"`
	Mark                   float64 `json:"mark"`
}
