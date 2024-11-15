package entity

import "gopkg.in/guregu/null.v4"

// ShipOrderStaging 发货台
type ShipOrderStaging struct {
	OrderDetailVOList []struct {
		ProductSkuId                int64       `json:"productSkuId"`                // 货品 skuId
		ProductSkuImgUrlList        []string    `json:"productSkuImgUrlList"`        // 货品 SKU 图片 URL 列表
		Color                       string      `json:"color"`                       // 颜色
		Size                        string      `json:"size"`                        // 尺码
		SkuDeliveryQuantityMaxLimit int         `json:"skuDeliveryQuantityMaxLimit"` // 发货数量限制最大值
		ProductOriginalSkuId        int64       `json:"productOriginalSkuId"`        // 原始 skuId
		ProductSkuPurchaseQuantity  int         `json:"productSkuPurchaseQuantity"`  // 货品 sku 下单数量
		PersonalText                null.String `json:"personalText"`                // 定制化内容
	} `json:"orderDetailVOList"` // 子订单详情信息
	SubPurchaseOrderBasicVO struct {
		SupplierId                         int64          `json:"supplierId"`
		IsCustomProduct                    bool           `json:"isCustomProduct"`                  // 是否为定制商品
		ExpectLatestArrivalTimeOrDefault   int            `json:"expectLatestArrivalTimeOrDefault"` // 要求最晚到达时间带默认值（时间戳 单位：毫秒）
		ProductSkcPicture                  string         `json:"productSkcPicture"`
		ProductName                        string         `json:"productName"`
		IsFirst                            bool           `json:"isFirst"`
		PurchaseStockType                  int            `json:"purchaseStockType"`
		DeliverUpcomingDelayTimeMillis     int64          `json:"deliverUpcomingDelayTimeMillis"`
		IsClothCategory                    bool           `json:"isClothCategory"`
		ProductSkcId                       int64          `json:"productSkcId"`
		SettlementType                     int            `json:"settlementType"`
		SkcExtCode                         string         `json:"skcExtCode"`
		DeliverDisplayCountdownMillis      int64          `json:"deliverDisplayCountdownMillis"`
		UrgencyType                        int            `json:"urgencyType"`
		SubWarehouseId                     int64          `json:"subWarehouseId"`
		ProductInventoryRegion             int            `json:"productInventoryRegion"`
		ExpectLatestDeliverTimeOrDefault   int            `json:"expectLatestDeliverTimeOrDefault"` // 要求最晚发货时间带默认值（时间戳 单位：毫秒）
		ArrivalUpcomingDelayTimeMillis     int64          `json:"arrivalUpcomingDelayTimeMillis"`
		ReceiveAddressInfo                 ReceiveAddress `json:"receiveAddressInfo"`
		AutoRemoveFromDeliveryPlatformTime int64          `json:"autoRemoveFromDeliveryPlatformTime"`
		ArrivalDisplayCountdownMillis      int64          `json:"arrivalDisplayCountdownMillis"`
		FragileTag                         bool           `json:"fragileTag"`
		PurchaseQuantity                   int            `json:"purchaseQuantity"`
		SubWarehouseName                   string         `json:"subWarehouseName"`
		PurchaseTime                       int64          `json:"purchaseTime"`
		SubPurchaseOrderSn                 string         `json:"subPurchaseOrderSn"`
	} `json:"subPurchaseOrderBasicVO"` // 子订单基本信息
}
