package entity

// ShipOrderPackingMatchResult 装箱发货校验结果
type ShipOrderPackingMatchResult struct {
	DeliveryOrderSnNotPrintBox     []string `json:"deliveryOrderSnNotPrintBox"`     // 未打印打包标签的发货单列表
	ShouldAddDeliveryOrderInfoList []string `json:"shouldAddDeliveryOrderInfoList"` // 需要勾选的相同发货地址的发货单列表（最多展示50个）
	SkuSumWeight                   int      `json:"skuSumWeight"`                   // 勾选的发货单对应SKU总重量（商品货品侧SKU重） 单位克
}
