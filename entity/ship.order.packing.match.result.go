package entity

// ShipOrderPackingMatchResult 装箱发货校验结果
type ShipOrderPackingMatchResult struct {
	ShouldAddDeliveryOrderTotal    int    `json:"shouldAddDeliveryOrderTotal"` // 需要勾选的相同发货地址的发货单的总个数
	TargetReceiveAddress           string `json:"targetReceiveAddress"`
	ShouldAddDeliveryOrderInfoList []struct {
		ReceiveSkcNum               int    `json:"receiveSkcNum"`
		ExpressPackageNum           int    `json:"expressPackageNum"`
		LatestFeedbackStatus        int    `json:"latestFeedbackStatus"`
		ExpressDeliverySn           string `json:"expressDeliverySn"`
		DeliveryOrderCancelLeftTime int64  `json:"deliveryOrderCancelLeftTime"`
		DeliveryAddressID           int    `json:"deliveryAddressId"`
		ExpressWeightFeedbackStatus int    `json:"expressWeightFeedbackStatus"`
		ExpressRejectStatus         int    `json:"expressRejectStatus"`
		PackageReceiveInfoVOList    []struct {
			ReceiveTime int64  `json:"receiveTime"`
			PackageSn   string `json:"packageSn"`
		} `json:"packageReceiveInfoVOList"`
		TaxWarehouseApplyOperateType int    `json:"taxWarehouseApplyOperateType"` // 入保税仓申请操作类型 0-不可操作 1-可申请 2-可查看
		ProductSkcID                 int64  `json:"productSkcId"`
		DeliveryContactAreaNo        string `json:"deliveryContactAreaNo"`
		SkcExtCode                   string `json:"skcExtCode"`
		InboundTime                  int64  `json:"inboundTime"`
		SubWarehouseID               int64  `json:"subWarehouseId"`
		PackageList                  []struct {
			SkcNum    int    `json:"skcNum"`
			PackageSn string `json:"packageSn"`
		} `json:"packageList"`
		InventoryRegion             int    `json:"inventoryRegion"`
		DeliverPackageNum           int    `json:"deliverPackageNum"`
		SubPurchaseOrderSn          string `json:"subPurchaseOrderSn"`
		DriverName                  string `json:"driverName"`
		ExpressCompanyID            int    `json:"expressCompanyId"`
		DefectiveSkcNum             int    `json:"defectiveSkcNum"`
		Status                      int    `json:"status"`
		ExpectPickUpGoodsTime       int    `json:"expectPickUpGoodsTime"`
		PredictTotalPackageWeight   int    `json:"predictTotalPackageWeight"`
		SupplierID                  int    `json:"supplierId"`
		IsDisplayCourier            bool   `json:"isDisplayCourier"`
		DeliveryMethod              int    `json:"deliveryMethod"`
		IsCustomProduct             bool   `json:"isCustomProduct"`
		ExpressWeightFeedbackTip    string `json:"expressWeightFeedbackTip"`
		ExceptionFeedBackTotalCount int    `json:"exceptionFeedBackTotalCount"`
		OtherDeliveryPackageNum     int    `json:"otherDeliveryPackageNum"`
		PurchaseStockType           int    `json:"purchaseStockType"`
		IfCanOperateDeliver         bool   `json:"ifCanOperateDeliver"`
		ReceivePackageNum           int    `json:"receivePackageNum"`
		IsPrintBoxMark              bool   `json:"isPrintBoxMark"`
		DeliveryContactNumber       string `json:"deliveryContactNumber"`
		ExpressCompany              string `json:"expressCompany"`
		IsClothCategory             bool   `json:"isClothCategory"`
		DeliveryOrderSn             string `json:"deliveryOrderSn"`
		DeliverTime                 int64  `json:"deliverTime"`
		UrgencyType                 int    `json:"urgencyType"`
		ExpressBatchSn              string `json:"expressBatchSn"`
		ReceiveAddressInfo          struct {
			DistrictCode  int    `json:"districtCode"`
			CityName      string `json:"cityName"`
			DistrictName  string `json:"districtName"`
			Phone         string `json:"phone"`
			ProvinceCode  int    `json:"provinceCode"`
			CityCode      int    `json:"cityCode"`
			ReceiverName  string `json:"receiverName"`
			DetailAddress string `json:"detailAddress"`
			ProvinceName  string `json:"provinceName"`
		} `json:"receiveAddressInfo"`
		PlateNumber       string `json:"plateNumber"`
		ReceiveTime       int64  `json:"receiveTime"`
		PackageDetailList []struct {
			ProductSkuID         int64  `json:"productSkuId"`
			ProductOriginalSkuId int64  `json:"productOriginalSkuId"`
			PersonalText         string `json:"personalText"`
			SkuNum               int    `json:"skuNum"`
		} `json:"packageDetailList"`
		DeliveryAddressInfo struct {
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
			ID            int    `json:"id"`
			ProvinceName  string `json:"provinceName"`
			WarehouseType int    `json:"warehouseType"`
		} `json:"deliveryAddressInfo"` // 需要勾选的相同发货地址的目标发货地址详情
		SubPurchaseOrderBasicVO struct {
			SupplierID                       int    `json:"supplierId"`
			IsCustomProduct                  bool   `json:"isCustomProduct"`
			ExpectLatestArrivalTimeOrDefault int    `json:"expectLatestArrivalTimeOrDefault"`
			ProductSkcPicture                string `json:"productSkcPicture"`
			ProductName                      string `json:"productName"`
			IsFirst                          bool   `json:"isFirst"`
			PurchaseStockType                int    `json:"purchaseStockType"`
			DeliverUpcomingDelayTimeMillis   int    `json:"deliverUpcomingDelayTimeMillis"`
			IsClothCategory                  bool   `json:"isClothCategory"`
			ProductSkcID                     int    `json:"productSkcId"`
			SettlementType                   int    `json:"settlementType"`
			SkcExtCode                       string `json:"skcExtCode"`
			DeliverDisplayCountdownMillis    int    `json:"deliverDisplayCountdownMillis"`
			UrgencyType                      int    `json:"urgencyType"`
			SubWarehouseID                   int    `json:"subWarehouseId"`
			ProductInventoryRegion           int    `json:"productInventoryRegion"`
			ExpectLatestDeliverTimeOrDefault int    `json:"expectLatestDeliverTimeOrDefault"`
			ReceiveAddressInfo               struct {
				DistrictCode  int    `json:"districtCode"`
				CityName      string `json:"cityName"`
				DistrictName  string `json:"districtName"`
				Phone         string `json:"phone"`
				ProvinceCode  int    `json:"provinceCode"`
				CityCode      int    `json:"cityCode"`
				ReceiverName  string `json:"receiverName"`
				DetailAddress string `json:"detailAddress"`
				ProvinceName  string `json:"provinceName"`
			} `json:"receiveAddressInfo"` // 需要勾选的相同收货地址的目标收货地址详情
			ArrivalUpcomingDelayTimeMillis     int64  `json:"arrivalUpcomingDelayTimeMillis"`
			AutoRemoveFromDeliveryPlatformTime int64  `json:"autoRemoveFromDeliveryPlatformTime"`
			ArrivalDisplayCountdownMillis      int64  `json:"arrivalDisplayCountdownMillis"`
			FragileTag                         bool   `json:"fragileTag"`
			PurchaseQuantity                   int    `json:"purchaseQuantity"`
			SubWarehouseName                   string `json:"subWarehouseName"`
			SubPurchaseOrderSn                 string `json:"subPurchaseOrderSn"`
			PurchaseTime                       int    `json:"purchaseTime"`
		} `json:"subPurchaseOrderBasicVO"`
		SubWarehouseName        string `json:"subWarehouseName"`
		PurchaseTime            int    `json:"purchaseTime"`
		SkcPurchaseNum          int    `json:"skcPurchaseNum"`
		DeliverSkcNum           int    `json:"deliverSkcNum"`
		DeliveryOrderCreateTime int64  `json:"deliveryOrderCreateTime"`
	} `json:"shouldAddDeliveryOrderInfoList"` // 需要勾选的相同发货地址的发货单列表（最多展示 50 个）
	AbleIgnorePlatformExpressForeMergeDelivery bool     `json:"ableIgnorePlatformExpressForeMergeDelivery"`
	TargetDeliveryAddress                      string   `json:"targetDeliveryAddress"`
	DeliveryOrderSnNotPrintBox                 []string `json:"deliveryOrderSnNotPrintBox"` // 未打印打包标签的发货单列表
	SkuSumWeight                               int      `json:"skuSumWeight"`               // 勾选的发货单对应SKU总重量（商品货品侧SKU重） 单位克
	AbleUsePlatformExpress                     bool     `json:"ableUsePlatformExpress"`     // 是否可以使用平台推荐物流服务商 true-表示商家无欠费可使用平台推荐物流
}
