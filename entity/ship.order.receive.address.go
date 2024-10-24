package entity

// ShipOrderReceiveAddress 发货单收货地址
type ShipOrderReceiveAddress struct {
	SubWarehouseId         int64          `json:"subWarehouseId"`         // 子仓 ID
	ReceiveAddressInfo     ReceiveAddress `json:"receiveAddressInfo"`     // 收货地址信息
	SubPurchaseOrderSnList []string       `json:"subPurchaseOrderSnList"` // 子采购单号列表
}
