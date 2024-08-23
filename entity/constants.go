package entity

// 发货地址类型
const (
	DeliveryAddressTypeChineseMainland = 1 // 内地
	DeliveryAddressTypeChineseHongKong = 2 // 香港
)

// 结算类型
const (
	SettlementTypeNotVMI = 0 // 非 VMI(采购)
	SettlementTypeVMI    = 1 // VMI(备货)
)

// 备货类型
const (
	PurchaseStockTypeNormal = 0 // 普通备货
	PurchaseStockTypeJIT    = 1 // JIT 备货
)

// 发货单紧急发货类型
const (
	ShipOrderTypeNormal  = 0 // 普通
	ShipOrderTypeUrgency = 1 // 加急
)

// 备货单状态
// 0-待接单；
// 1-已接单，待发货；
// 2-已送货；
// 3-已收货；
// 4-已拒收；
// 5-已验收，全部退回；
// 6-已验收；
// 7-已入库；
// 8-作废；
// 9-已超时；
// 10-已取消
const (
	PurchaseOrderStatusWaitingMerchantReceive = 0  // 待接单
	PurchaseOrderStatusMerchantReceived       = 1  // 已接单/待发货
	PurchaseOrderStatusMerchantSend           = 2  // 已送货
	PurchaseOrderStatusPlatformReceived       = 3  // 已收货
	PurchaseOrderStatusPlatformRejected       = 4  // 已拒收
	PurchaseOrderStatusPlatformReturned       = 5  // 已验收/全部退回
	PurchaseOrderStatusPlatformApproved       = 6  // 已验收
	PurchaseOrderStatusPlatformPutInStorage   = 7  // 已入库
	PurchaseOrderStatusDiscard                = 8  // 作废
	PurchaseOrderStatusTimeout                = 9  // 已超时
	PurchaseOrderStatusCancel                 = 10 // 已取消
)

// 备货单来源
const (
	PurchaseOrderSourceOperationalStaff = 0    // 运营下单
	PurchaseOrderSourceSeller           = 1    // 卖家下单
	PurchaseOrderSourcePlatform         = 9999 // 平台下单
)

// 发货方式
const (
	DeliveryMethodSelf                   = 1 // 自行配送
	DeliveryMethodPlatformRecommendation = 2 // 平台推荐服务商
	DeliveryMethodThirdParty             = 3 // 自行委托第三方物流
)
