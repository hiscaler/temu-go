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
		IsCustomProduct                    bool           `json:"isCustomProduct"`                    // 是否为定制商品
		ExpectLatestArrivalTimeOrDefault   int64          `json:"expectLatestArrivalTimeOrDefault"`   // 要求最晚到达时间带默认值（时间戳 单位：毫秒）
		ProductSkcPicture                  string         `json:"productSkcPicture"`                  // 货品图片
		ProductName                        string         `json:"productName"`                        // 商品名
		IsFirst                            bool           `json:"isFirst"`                            // 是否首单
		PurchaseStockType                  int            `json:"purchaseStockType"`                  // 备货类型 0-普通备货 1-jit备货
		DeliverUpcomingDelayTimeMillis     int64          `json:"deliverUpcomingDelayTimeMillis"`     // 剩余发货时间不足XX，则统计为即将逾期，前端展示标红 单位：毫秒 默认12 * 3600 * 1000
		IsClothCategory                    bool           `json:"isClothCategory"`                    // 是否服饰类目
		ProductSkcId                       int64          `json:"productSkcId"`                       // skcId
		SettlementType                     int            `json:"settlementType"`                     // 结算类型 0-非vmi 1-vmi
		SkcExtCode                         string         `json:"skcExtCode"`                         // 货号
		DeliverDisplayCountdownMillis      int64          `json:"deliverDisplayCountdownMillis"`      // 剩余发货时间不足XX，则前端开始读秒 单位：毫秒 默认1小时
		UrgencyType                        int            `json:"urgencyType"`                        // 是否是紧急发货单，0-普通 1-急采
		SubWarehouseId                     int64          `json:"subWarehouseId"`                     // 子仓 id
		ProductInventoryRegion             int            `json:"productInventoryRegion"`             // 备货类型
		ExpectLatestDeliverTimeOrDefault   int64          `json:"expectLatestDeliverTimeOrDefault"`   // 要求最晚发货时间带默认值（时间戳 单位：毫秒）
		ArrivalUpcomingDelayTimeMillis     int64          `json:"arrivalUpcomingDelayTimeMillis"`     // 剩余到货时间不足XX，则统计为即将逾期，前端展示标红 单位：毫秒 默认6 * 3600 * 1000
		ReceiveAddressInfo                 ReceiveAddress `json:"receiveAddressInfo"`                 // 收货仓详细地址
		AutoRemoveFromDeliveryPlatformTime int64          `json:"autoRemoveFromDeliveryPlatformTime"` // 自动移出发货台倒计时时间,毫秒
		ArrivalDisplayCountdownMillis      int64          `json:"arrivalDisplayCountdownMillis"`      // 剩余到货时间不足XX，则前端开始读秒 单位：毫秒 默认1小时
		FragileTag                         bool           `json:"fragileTag"`                         // 易碎品打标
		PurchaseQuantity                   int            `json:"purchaseQuantity"`                   // 下单数量
		SubWarehouseName                   string         `json:"subWarehouseName"`                   // 子仓名称
		PurchaseTime                       int64          `json:"purchaseTime"`                       // 下单时间：毫秒
		SubPurchaseOrderSn                 string         `json:"subPurchaseOrderSn"`                 // 采购子单号
	} `json:"subPurchaseOrderBasicVO"` // 子订单基本信息
	OrderType null.Int `json:"orderType"` // 自处理后的发货台数据类型（1：普通、2：JIT、3：定制）
}
