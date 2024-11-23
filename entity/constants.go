package entity

const (
	MaxPageSize = 100 // 每页最多数量
)

// 布尔数据表示值
const (
	FalseNumber = 0
	TrueNumber  = 1
)

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

// 扩展备货类型
const (
	StockTypeNormal     = 1 // 普通备货单
	StockTypeJIT        = 2 // JIT 备货单
	StockTypeCustomized = 3 // 定制备货单
)

// 紧急发货单类型
const (
	UrgencyTypeNormal  = 0 // 普通
	UrgencyTypeUrgency = 1 // 加急
)

// 备货单状态
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
	DeliveryMethodNone                   = 0 // 无
	DeliveryMethodSelf                   = 1 // 自行配送
	DeliveryMethodPlatformRecommendation = 2 // 平台推荐服务商
	DeliveryMethodThirdParty             = 3 // 自行委托第三方物流
)

// 发货单状态
const (
	ShipOrderStatusWaitingPacking          = 0 // 待装箱发货
	ShipOrderStatusWaitingWarehouseReceive = 1 // 待仓库收货
	ShipOrderStatusReceived                = 2 // 已收货
	ShipOrderStatusInStorage               = 3 // 已入库
	ShipOrderStatusReturned                = 4 // 已退货
	ShipOrderStatusCanceled                = 5 // 已取消
	ShipOrderStatusPartialReceive          = 6 // 部分收货
)
