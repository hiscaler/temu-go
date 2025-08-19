package entity

import "gopkg.in/guregu/null.v4"

// ShipOrderStaging 发货台
type ShipOrderStaging struct {
	OrderDetailVOList       []ShipOrderStagingOrderDetail `json:"orderDetailVOList"`       // 子订单详情信息
	SubPurchaseOrderBasicVO SubPurchaseOrderBasic         `json:"subPurchaseOrderBasicVO"` // 子订单基本信息
	OrderType               null.Int                      `json:"orderType"`               // 自处理后的发货台数据类型（1：普通、2：紧急、3：定制）
}

// ShipOrderStagingOrderDetail 子订单详情信息
type ShipOrderStagingOrderDetail struct {
	ProductSkuId                int64       `json:"productSkuId"`                // 货品 skuId
	ProductSkuImgUrlList        []string    `json:"productSkuImgUrlList"`        // 货品 SKU 图片 URL 列表
	Color                       string      `json:"color"`                       // 颜色
	Size                        string      `json:"size"`                        // 尺码
	SkuDeliveryQuantityMaxLimit int         `json:"skuDeliveryQuantityMaxLimit"` // 发货数量限制最大值
	ProductOriginalSkuId        int64       `json:"productOriginalSkuId"`        // 原始 skuId
	ProductSkuPurchaseQuantity  int         `json:"productSkuPurchaseQuantity"`  // 货品 sku 下单数量
	PersonalText                null.String `json:"personalText"`                // 定制化内容
}
