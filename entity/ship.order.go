package entity

import "gopkg.in/guregu/null.v4"

// ShipOrder 发货单
type ShipOrder struct {
	ReceiveSkcNum                int      `json:"receiveSkcNum"`
	ExpressPackageNum            int      `json:"expressPackageNum"`
	LatestFeedbackStatus         int      `json:"latestFeedbackStatus"`
	ExpectLatestPickTime         null.Int `json:"expectLatestPickTime"`
	DeliveryOrderCancelLeftTime  null.Int `json:"deliveryOrderCancelLeftTime"`
	ExpressDeliverySn            string   `json:"expressDeliverySn"`
	DeliveryAddressId            null.Int `json:"deliveryAddressId"`
	ExpressWeightFeedbackStatus  int      `json:"expressWeightFeedbackStatus"`
	ExpressRejectStatus          null.Int `json:"expressRejectStatus"`
	PackageReceiveInfoVOList     any      `json:"packageReceiveInfoVOList"`
	TaxWarehouseApplyOperateType int      `json:"taxWarehouseApplyOperateType"`
	ProductSkcId                 int64    `json:"productSkcId"`
	SkcExtCode                   string   `json:"skcExtCode"`
	InboundTime                  null.Int `json:"inboundTime"`
	SubWarehouseId               int64    `json:"subWarehouseId"`
	PackageList                  []struct {
		SkcNum    int    `json:"skcNum"`
		PackageSn string `json:"packageSn"`
	} `json:"packageList"`
	InventoryRegion             int            `json:"inventoryRegion"`
	DeliverPackageNum           int            `json:"deliverPackageNum"`
	DriverName                  string         `json:"driverName"`
	SubPurchaseOrderSn          string         `json:"subPurchaseOrderSn"`
	ExpressCompanyId            int64          `json:"expressCompanyId"`
	DefectiveSkcNum             int            `json:"defectiveSkcNum"`
	Status                      int            `json:"status"`
	ExpectPickUpGoodsTime       int            `json:"expectPickUpGoodsTime"`
	PredictTotalPackageWeight   int64          `json:"predictTotalPackageWeight"`
	SupplierId                  int64          `json:"supplierId"`
	IsDisplayCourier            bool           `json:"isDisplayCourier"`
	IsCustomProduct             bool           `json:"isCustomProduct"`
	DeliveryMethod              int            `json:"deliveryMethod"`
	ExpressWeightFeedbackTip    any            `json:"expressWeightFeedbackTip"`
	ExceptionFeedBackTotalCount null.Int       `json:"exceptionFeedBackTotalCount"`
	OtherDeliveryPackageNum     int            `json:"otherDeliveryPackageNum"`
	PurchaseStockType           int            `json:"purchaseStockType"`
	IfCanOperateDeliver         bool           `json:"ifCanOperateDeliver"`
	ReceivePackageNum           int            `json:"receivePackageNum"`
	IsPrintBoxMark              bool           `json:"isPrintBoxMark"`
	ExpressCompany              string         `json:"expressCompany"`
	IsClothCategory             bool           `json:"isClothCategory"`
	DeliveryOrderSn             string         `json:"deliveryOrderSn"`
	DeliverTime                 null.Int       `json:"deliverTime"`
	UrgencyType                 int            `json:"urgencyType"`
	ExpressBatchSn              string         `json:"expressBatchSn"`
	ReceiveAddressInfo          ReceiveAddress `json:"receiveAddressInfo"`
	PlateNumber                 string         `json:"plateNumber"`
	ReceiveTime                 null.Int       `json:"receiveTime"`
	PackageDetailList           []struct {
		ProductSkuId         int64       `json:"productSkuId"`
		ProductOriginalSkuId null.Int    `json:"productOriginalSkuId"`
		PersonalText         null.String `json:"personalText"`
		SkuNum               int         `json:"skuNum"`
	} `json:"packageDetailList"`
	SubPurchaseOrderBasicVO struct {
		SupplierId         int64       `json:"supplierId"`
		IsCustomProduct    bool        `json:"isCustomProduct"`
		ProductSkcPicture  string      `json:"productSkcPicture"`
		IsFirst            bool        `json:"isFirst"`
		PurchaseStockType  int         `json:"purchaseStockType"`
		IsClothCategory    bool        `json:"isClothCategory"`
		ProductSkcId       int64       `json:"productSkcId"`
		SettlementType     int         `json:"settlementType"`
		SkcExtCode         string      `json:"skcExtCode"`
		SubWarehouseId     null.Int    `json:"subWarehouseId"`
		UrgencyType        int         `json:"urgencyType"`
		FragileTag         bool        `json:"fragileTag"`
		PurchaseQuantity   int         `json:"purchaseQuantity"`
		SubWarehouseName   null.String `json:"subWarehouseName"`
		PurchaseTime       int64       `json:"purchaseTime"`
		SubPurchaseOrderSn string      `json:"subPurchaseOrderSn"`
	} `json:"subPurchaseOrderBasicVO"`
	SubWarehouseName        string `json:"subWarehouseName"`
	PurchaseTime            int64  `json:"purchaseTime"`
	SkcPurchaseNum          int    `json:"skcPurchaseNum"`
	DeliverSkcNum           int    `json:"deliverSkcNum"`
	DeliveryOrderCreateTime int64  `json:"deliveryOrderCreateTime"`
}
