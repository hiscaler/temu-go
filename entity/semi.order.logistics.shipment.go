package entity

// ShipmentResult 包含物流发货信息的结果
type ShipmentResult struct {
	Result struct {
		ShipmentInfoDTO []ShipmentInfoDTO `json:"shipmentInfoDTO"` // 发货信息列表
	} `json:"result"`
}

// ShipmentInfo 表示单个物流发货信息
type ShipmentInfoDTO struct {
	CarrierId            int64                `json:"carrierId"`                  // 物流公司ID
	CarrierName          string               `json:"carrierName"`                // 物流公司名称
	TrackingNumber       string               `json:"trackingNumber"`             // 运单号
	SkuId                int64                `json:"skuId"`                      // 商品skuId
	Quantity             int                  `json:"quantity"`                   // 商品skuId对应发货数量
	PackageSn            string               `json:"packageSn"`                  // 包裹号
	PackageDeliveryType  int                  `json:"packageDeliveryType"`        // 发货包裹履约类型：1-导入运单发货 2-在线下单发货 3-合作对接仓导入运单发货 4-合作对接仓在线下单发货
	CooperativeWarehouse CooperativeWarehouse `json:"cooperativeWarehouseDTO"`    // 合作对接仓信息，仅在packageDeliveryType=3或4时返回
	TrackingWarningLabel int                  `json:"trackingWarningLabel"`       // 运单物流提醒标签：0-无问题 1-查无轨迹 2-疑似有误 3-收货地址不一致 4-未揽收
	SubPackageShipments  []SubPackageShipment `json:"subPackageShipmentInfoList"` // 附属包裹列表，当为单sku拆单发货场景时，后续增加补充的运单信息将作为附属包裹展示
}

// CooperativeWarehouse 表示合作对接仓信息
type CooperativeWarehouse struct {
	WarehouseProviderCode      string `json:"warehouseProviderCode"`      // 合作对接仓服务商编码
	WarehouseProviderBrandName string `json:"warehouseProviderBrandName"` // 合作对接仓服务商名字
	WarehouseCode              string `json:"warehouseCode"`              // 合作对接仓编码
	WarehouseName              string `json:"warehouseName"`              // 合作对接仓名字
}

// SubPackageShipment 表示附属包裹发货信息
type SubPackageShipment struct {
	TrackingNumber      string `json:"trackingNumber"`      // 运单号
	CarrierId           int64  `json:"carrierId"`           // 物流公司ID
	CarrierName         string `json:"carrierName"`         // 物流公司名称
	PackageSn           string `json:"packageSn"`           // 包裹号
	PackageDeliveryType int    `json:"packageDeliveryType"` // 发货包裹履约类型：1-导入运单发货 2-在线下单发货 3-合作对接仓导入运单发货 4-合作对接仓在线下单发货
}
