package entity

type SemiOnlineOrderPlatformLogisticsUnshippedPackage struct {
	CarrierId   int64  `json:"carrierId"`   // 物流渠道 ID
	CarrierName string `json:"carrierName"` // 物流渠道名称
	// 对应枚举如下：
	//	主包裹：MAIN
	//	子包裹：SPLIT_LARGE_ITEM
	//	备注：
	//	1、当包裹下call是单件SKU多包裹的形式时，会将【sendSubRequestList】入参的相关包裹视为关联这个包裹的子包裹。
	//	2、默认所有非单sku多包裹场景的包裹都是主包裹。
	SubPackageType   string `json:"subPackageType"` // 包裹主子类型
	MainPackageSn    string `json:"mainPackageSn"`  // 该包裹关联的主包裹。当包裹为主包裹时返回自身，当包裹为子包裹时返回关联的主包裹
	PackageSn        string `json:"packageSn"`      // 包裹号
	TrackingNumber   string `json:"trackingNumber"` // 运单号
	SubPackageSnList []struct {
		ShippableOrders []struct {
			ParentOrderSn string `json:"parentOrderSn"` // 父单号
			OrderSn       string `json:"orderSn"`       // 子单号
			Quantity      int    `json:"quantity"`      // 应履约件数
		} `json:"shippableOrders"` // 包裹内能发运的O单详情
		CanceledOrders []struct {
			ParentOrderSn string `json:"parentOrderSn"` // 父单号
			OrderSn       string `json:"orderSn"`       // 子单号
			Quantity      int    `json:"quantity"`      // 应履约件数
		} `json:"canceledOrders"` // 包裹内已被取消的O单详情
	} `json:"subPackageSnList"` // 该包裹下的子包裹列表。当包裹为主包裹时返回关联的子包裹，当包裹为子包裹时返回为空。
	PackageDetail struct {
		CanceledOrders  []any `json:"canceledOrders"`
		ShippableOrders []struct {
			ParentOrderSn string `json:"parentOrderSn"`
			OrderSn       string `json:"orderSn"`
			Quantity      int    `json:"quantity"`
		} `json:"shippableOrders"`
	} `json:"packageDetail"` // 包裹详情
}
