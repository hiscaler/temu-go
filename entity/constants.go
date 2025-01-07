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

// 扩展备货单、发货单类型
const (
	OrderTypeNormal     = 1 // 普通备货单
	OrderTypeJIT        = 2 // JIT 备货单
	OrderTypeCustomized = 3 // 定制备货单
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

// 备货单逾期状态
const (
	PurchaseOrderDeliveryWillBeDelay = 101 // 发货即将逾期
	PurchaseOrderDeliveryDelay       = 102 // 发货已逾期
	PurchaseOrderArrivalWillBeDelay  = 201 // 到货即将逾期
	PurchaseOrderArrivalDelay        = 202 // 到货已逾期
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

// 商品选品状态
const (
	GoodsSelectionStatusWaitingForFirstOrder    = 10 // 待下首单
	GoodsSelectionStatusHasBeenPlacedFirstOrder = 11 // 已下首单
	GoodsSelectionStatusJoinedSite              = 12 // 已加入站点
	GoodsSelectionStatusOffline                 = 13 // 已下架
)

// 商品售罄状态
const (
	GoodsStockStatusAbundant      = 0 // 库存充足
	GoodsStockStatusWillBeSoldOut = 1 // 即将售罄
	GoodsStockStatusSoldOut       = 2 // 已经售罄
)

const (
	LanguageZhCn = "zh" // 中文
	LanguageEn   = "en" // 英文
)

// 物流面单文件类型
const (
	LogisticsShipmentDocumentPdfFile   = "SHIPPING_LABEL_PDF"
	LogisticsShipmentDocumentImageFile = "SHIPPING_LABEL_IMAGE"
)

// 数量增减模式
const (
	QuantityChangeModeInDecrease = 1 // 增减变量
	QuantityChangeModeReplace    = 2 // 覆盖变更
)

// https://seller.kuajingmaihuo.com/sop/view/231998342274104483#6mTvhA
const (
	ChinaRegion         = "CN" // 中国区
	AmericanRegion      = "US" // 美区
	EuropeanUnionRegion = "EU" // 欧盟区
)

// https://seller.kuajingmaihuo.com/sop/view/231998342274104483#d78RUG
// 站点 Id
const (
	AmericanSiteId      = 100 // 美国站
	CanadaSiteId        = 101 // 加拿大站
	UnitedKingdomSiteId = 102 // 英国站
	AustraliaSiteId     = 103 // 澳大利亚站
)

var AmericanSiteIds = []int{100, 101, 103, 104, 118, 110, 187}                                                                                                               // 美区站点 ID
var EuropeanUnionSiteIds = []int{105, 106, 107, 109, 102, 112, 137, 138, 111, 108, 142, 113, 143, 140, 139, 145, 116, 146, 141, 115, 144, 150, 148, 147, 149, 151, 117, 152} // 欧盟区站点 ID
var SiteIds = []int{
	100, 101, 103, 104, 118, 110, 187,
	105, 106, 107, 109, 102, 112, 137, 138, 111, 108, 142, 113, 143, 140, 139, 145, 116, 146, 141, 115, 144, 150, 148, 147, 149, 151, 117, 152,
} // 所有区站点 ID
