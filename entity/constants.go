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
	OrderTypeUrgent     = 2 // 紧急备货单
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

// 半托管订单排序
const (
	SemiOrderOrderByCreateTime = "createTime"
	SemiOrderOrderByUpdateTime = "updateTime"
)

// 子订单履约类型
const (
	SemiOrderFulfillmentTypeBySeller               = "fulfillBySeller"               // 只返回卖家履约子订单列表
	SemiOrderFulfillmentTypeByCooperativeWarehouse = "fulfillByCooperativeWarehouse" // 只返回合作仓履约子订单列表
)

// 半托管 PO 单标签
const (
	SemiParentOrderLabelSoonToBeOverdue           = "soon_to_be_overdue"           // 即将逾期
	SemiParentOrderLabelPastDue                   = "past_due"                     // 已逾期
	SemiParentOrderLabelPendingBuyerCancellation  = "pending_buyer_cancellation"   // 买家取消待确认订单
	SemiParentOrderLabelPendingBuyerAddressChange = "pending_buyer_address_change" // 买家改地址待确认订单
)

// 半托物流发货类型
const (
	SemiShippingTypeSingle = 0 // 单个运单发货
	SemiShippingTypeSplit  = 1 // 拆成多个运单发货
	SemiShippingTypeMerge  = 2 // 合并发货
)

// TEMU 半托管订单状态（orderStatus、parentOrderStatus）
// https://seller.kuajingmaihuo.com/sop/view/867739977041685428#E10GkB
const (
	SemiOrderStatusAll              = 0  // 全部
	SemiOrderStatusPending          = 1  // 订单挂起中，用户支付后PO单(parentOrder)进入pending状态，订单暂时被挂起，用户支付后存在一段时间下单冷静期，包含风控处置等订单处理时间
	SemiOrderStatusUnShipping       = 2  // 订单待发货，待发货状态下用户依然可以取消订单，需要及时监听订单取消状态及时接单
	SemiOrderStatusCanceled         = 3  // 订单已取消，用户已经取消订单
	SemiOrderStatusShipped          = 4  // 订单已发货，订单已经发货完成
	SemiOrderStatusReceipted        = 5  // 订单已签收
	SemiOrderStatusPartialCanceled  = 41 // 部分取消（本本订单）
	SemiOrderStatusPartialReceipted = 51 // 部分签收（本本订单）
)

// 半托定制模板类型
const (
	SemiOrderCustomizationTemplateTypeImage        = 1 // Only image
	SemiOrderCustomizationTemplateTypeText         = 2 // Only text
	SemiOrderCustomizationTemplateTypeTextAndImage = 3 // Text and image
)

// 半托定制类型
const (
	SemiOrderCustomizationCustomizedTypePureText        = 1 // pure text customization, no customized templates
	SemiOrderCustomizationCustomizedTypeGraphicsAndText = 2 // customized graphics and text, with customized templates available
)

// https://seller.kuajingmaihuo.com/sop/view/231998342274104483#6mTvhA
const (
	ChinaRegion         = "CN"      // 中国区
	AmericanRegion      = "US"      // 美区
	EuropeanUnionRegion = "EU"      // 欧盟区
	GlobalRegion        = "GLOBAL"  // 全球（除欧区、美国）
	Partner             = "PARTNER" // 合作伙伴
)

var RegionIds = []int{
	2, 3, 4, 5, 6, 7, 8, 9, 10,
	11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
	21, 22, 23, 24, 25, 26, 27, 28, 29, 30,
	31, 32, 33, 34, 35, 36, 37, 38, 39, 40,
	41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
	51, 52, 53, 54, 55, 56, 57, 58, 59, 60,
	61, 62, 63, 64, 65, 66, 67, 68, 69, 70,
	71, 72, 73, 74, 75, 76, 77, 78, 79, 80,
	81, 82, 83, 84, 85, 86, 87, 88, 89, 90,
	91, 92, 93, 94, 95, 96, 97, 98, 99, 100,
	101, 102, 103, 104, 105, 106, 107, 108, 109, 110,
	111, 112, 113, 114, 115, 116, 117, 118, 119, 120,
	121, 122, 123, 124, 125, 126, 127, 128, 129, 130,
	131, 132, 133, 134, 135, 136, 137, 138, 139, 140,
	141, 142, 143, 144, 145, 146, 147, 148, 149, 150,
	151, 152, 153, 154, 155, 156, 157, 158, 159, 160,
	161, 162, 163, 164, 165, 166, 167, 168, 169, 170,
	171, 172, 173, 174, 175, 176, 177, 178, 179, 180,
	181, 182, 183, 184, 185, 186, 187, 188, 189, 190,
	191, 192, 193, 194, 195, 196, 197, 198, 199, 200,
	201, 202, 203, 204, 205, 206, 207, 208, 209, 210,
	211, 212, 213, 214, 215, 216, 217, 218, 219, 220,
	221, 222, 223, 224, 225, 226,
}

// https://seller.kuajingmaihuo.com/sop/view/231998342274104483#d78RUG
// 站点 Id
const (
	AmericanSiteId      = 100 // 美国站
	CanadaSiteId        = 101 // 加拿大站
	UnitedKingdomSiteId = 102 // 英国站
	AustraliaSiteId     = 103 // 澳大利亚站
	SaudiArabiaSiteId   = 120 // 沙特站
	BelgiumSiteId       = 142 // 比利时站
	VietNamSiteId       = 187 // 越南站
)

var AmericanSiteIds = []int{100, 101, 103, 104, 118, 110, 187}                                                                                                               // 美区站点 ID
var EuropeanUnionSiteIds = []int{105, 106, 107, 109, 102, 112, 137, 138, 111, 108, 142, 113, 143, 140, 139, 145, 116, 146, 141, 115, 144, 150, 148, 147, 149, 151, 117, 152} // 欧盟区站点 ID
var SiteIds = []int{
	100, 101, 103, 104, 118, 110, 120, 187,
	105, 106, 107, 109, 102, 112, 137, 138, 111, 108, 142, 113, 143, 140, 139, 145, 116, 146, 141, 115, 144, 150, 148, 147, 149, 151, 117, 152,
} // 所有区站点 ID
