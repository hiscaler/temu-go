package entity

// ShipOrderPackingMatchResult 装箱发货校验结果
type ShipOrderPackingMatchResult struct {
	DeliveryOrderSnNotPrintBox                 []string                                      `json:"deliveryOrderSnNotPrintBox"`                 // 未打印打包标签的发货单列表
	ShouldAddDeliveryOrderInfoList             []ShipOrderPackingMatchShouldAddDeliveryOrder `json:"shouldAddDeliveryOrderInfoList"`             // 需要勾选的相同发货地址的发货单列表（最多展示 50 个）
	TargetReceiveAddress                       string                                        `json:"targetReceiveAddress"`                       // 需要勾选的相同收货地址的目标收货地址详情
	TargetDeliveryAddress                      string                                        `json:"targetDeliveryAddress"`                      // 需要勾选的相同发货地址的目标发货地址详情
	ShouldAddDeliveryOrderTotal                int                                           `json:"shouldAddDeliveryOrderTotal"`                // 需要勾选的相同发货地址的发货单的总个数
	AbleUsePlatformExpress                     bool                                          `json:"ableUsePlatformExpress"`                     // 是否可以使用平台推荐物流服务商 true-表示商家无欠费可使用平台推荐物流
	AbleIgnorePlatformExpressForeMergeDelivery bool                                          `json:"ableIgnorePlatformExpressForeMergeDelivery"` // 是否可以忽略平台物流强制合并发货:true/false/null, null-未命中判断灰度 true-发货单SKU总重量大于50KG,使用平台物流时可以忽略强制合并发货
	SkuSumWeight                               int64                                         `json:"skuSumWeight"`                               // 勾选的发货单对应SKU总重量（商品货品侧SKU重） 单位克
}

// ShipOrderPackingMatchShouldAddDeliveryOrder 装箱发货校验相同发货地址的发货单
type ShipOrderPackingMatchShouldAddDeliveryOrder struct {
	ReceiveSkcNum               int    `json:"receiveSkcNum"`
	ExpressPackageNum           int    `json:"expressPackageNum"`
	LatestFeedbackStatus        int    `json:"latestFeedbackStatus"`
	ExpressDeliverySn           string `json:"expressDeliverySn"`
	DeliveryOrderCancelLeftTime int64  `json:"deliveryOrderCancelLeftTime"`
	DeliveryAddressId           int64  `json:"deliveryAddressId"`
	ExpressWeightFeedbackStatus int    `json:"expressWeightFeedbackStatus"`
	ExpressRejectStatus         int    `json:"expressRejectStatus"`
	PackageReceiveInfoVOList    []struct {
		ReceiveTime int64  `json:"receiveTime"`
		PackageSn   string `json:"packageSn"`
	} `json:"packageReceiveInfoVOList"`
	TaxWarehouseApplyOperateType int    `json:"taxWarehouseApplyOperateType"` // 入保税仓申请操作类型 0-不可操作 1-可申请 2-可查看
	ProductSkcId                 int64  `json:"productSkcId"`
	DeliveryContactAreaNo        string `json:"deliveryContactAreaNo"`
	SkcExtCode                   string `json:"skcExtCode"`
	InboundTime                  int64  `json:"inboundTime"`
	SubWarehouseId               int64  `json:"subWarehouseId"`
	PackageList                  []struct {
		SkcNum    int    `json:"skcNum"`
		PackageSn string `json:"packageSn"`
	} `json:"packageList"`
	InventoryRegion             int            `json:"inventoryRegion"`
	DeliverPackageNum           int            `json:"deliverPackageNum"`
	SubPurchaseOrderSn          string         `json:"subPurchaseOrderSn"`
	DriverName                  string         `json:"driverName"`
	ExpressCompanyId            int64          `json:"expressCompanyId"`
	DefectiveSkcNum             int            `json:"defectiveSkcNum"`
	Status                      int            `json:"status"`
	ExpectPickUpGoodsTime       int64          `json:"expectPickUpGoodsTime"`
	PredictTotalPackageWeight   int            `json:"predictTotalPackageWeight"`
	SupplierId                  int            `json:"supplierId"`
	IsDisplayCourier            bool           `json:"isDisplayCourier"`
	DeliveryMethod              int            `json:"deliveryMethod"`
	IsCustomProduct             bool           `json:"isCustomProduct"`
	ExpressWeightFeedbackTip    string         `json:"expressWeightFeedbackTip"`
	ExceptionFeedBackTotalCount int            `json:"exceptionFeedBackTotalCount"`
	OtherDeliveryPackageNum     int            `json:"otherDeliveryPackageNum"`
	PurchaseStockType           int            `json:"purchaseStockType"`
	IfCanOperateDeliver         bool           `json:"ifCanOperateDeliver"`
	ReceivePackageNum           int            `json:"receivePackageNum"`
	IsPrintBoxMark              bool           `json:"isPrintBoxMark"`
	DeliveryContactNumber       string         `json:"deliveryContactNumber"`
	ExpressCompany              string         `json:"expressCompany"`
	IsClothCategory             bool           `json:"isClothCategory"`
	DeliveryOrderSn             string         `json:"deliveryOrderSn"`
	DeliverTime                 int64          `json:"deliverTime"`
	UrgencyType                 int            `json:"urgencyType"`
	ExpressBatchSn              string         `json:"expressBatchSn"`
	ReceiveAddressInfo          ReceiveAddress `json:"receiveAddressInfo"`
	PlateNumber                 string         `json:"plateNumber"`
	ReceiveTime                 int64          `json:"receiveTime"`
	PackageDetailList           []struct {
		ProductSkuId         int64  `json:"productSkuId"`
		ProductOriginalSkuId int64  `json:"productOriginalSkuId"`
		PersonalText         string `json:"personalText"`
		SkuNum               int    `json:"skuNum"`
	} `json:"packageDetailList"`
	DeliveryAddressInfo struct {
		ID            int    `json:"id"`
		TownName      string `json:"townName"`
		DistrictCode  int    `json:"districtCode"`
		DistrictName  string `json:"districtName"`
		TownCode      int    `json:"townCode"`
		ProvinceCode  int    `json:"provinceCode"`
		CityCode      int    `json:"cityCode"`
		AddressLabel  string `json:"addressLabel"`
		IsDefault     bool   `json:"isDefault"`
		AddressDetail string `json:"addressDetail"`
		CityName      string `json:"cityName"`
		ProvinceName  string `json:"provinceName"`
		WarehouseType int    `json:"warehouseType"`
	} `json:"deliveryAddressInfo"` // 需要勾选的相同发货地址的目标发货地址详情
	SubPurchaseOrderBasicVO SubPurchaseOrderBasic `json:"subPurchaseOrderBasicVO"`
	SubWarehouseName        string                `json:"subWarehouseName"`
	PurchaseTime            int64                 `json:"purchaseTime"`
	SkcPurchaseNum          int                   `json:"skcPurchaseNum"`
	DeliverSkcNum           int                   `json:"deliverSkcNum"`
	DeliveryOrderCreateTime int64                 `json:"deliveryOrderCreateTime"`
}
