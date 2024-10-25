package entity

import "gopkg.in/guregu/null.v4"

// BoxMarkInfo 箱唛
type BoxMarkInfo struct {
	VolumeType                  int             `json:"volumeType"`
	SupplierId                  int             `json:"supplierId"`
	DeliveryMethod              int             `json:"deliveryMethod"`
	IsCustomProduct             bool            `json:"isCustomProduct"`
	PackageIndex                int             `json:"packageIndex"`
	NonClothMainSpecVOList      []Specification `json:"nonClothMainSpecVOList"`
	ExpressDeliverySn           string          `json:"expressDeliverySn"`
	ProductName                 string          `json:"productName"`
	SubWarehouseEnglishName     string          `json:"subWarehouseEnglishName"`
	IsClothCat                  bool            `json:"isClothCat"`
	IsFirst                     bool            `json:"isFirst"`
	PurchaseStockType           int             `json:"purchaseStockType"`
	TotalPackageNum             int             `json:"totalPackageNum"`
	ExpressCompany              string          `json:"expressCompany"`
	ProductSkcId                int64           `json:"productSkcId"`
	NonClothSkuExtCode          string          `json:"nonClothSkuExtCode"`
	DeliveryOrderSn             string          `json:"deliveryOrderSn"`
	SettlementType              int             `json:"settlementType"`
	SupplierName                string          `json:"supplierName"`
	ProductSkuIdList            []int64         `json:"productSkuIdList"`
	SkcExtCode                  string          `json:"skcExtCode"`
	DeliverTime                 int             `json:"deliverTime"`
	UrgencyType                 int             `json:"urgencyType"`
	SubWarehouseId              int64           `json:"subWarehouseId"`
	ProductSkcName              null.String     `json:"productSkcName"`
	PackageSn                   string          `json:"packageSn"`
	ExpressEnglishCompany       string          `json:"expressEnglishCompany"`
	PackageSkcNum               int             `json:"packageSkcNum"`
	NonClothSecondarySpecVOList []Specification `json:"nonClothSecondarySpecVOList"`
	SubWarehouseName            string          `json:"subWarehouseName"`
	DriverName                  string          `json:"driverName"`
	SubPurchaseOrderSn          string          `json:"subPurchaseOrderSn"`
	DriverPhone                 null.String     `json:"driverPhone"`
	PurchaseTime                int64           `json:"purchaseTime"`
	GreyKeyHitMap               struct {
		KeyCnBG215946 bool `json:"key_cn_BG_215946"`
		KeyBG31572    bool `json:"key_BG_31572"`
		KeyCnBG226344 bool `json:"key_cn_BG_226344"`
	} `json:"greyKeyHitMap"`
	StorageAttrName any `json:"storageAttrName"`
	DeliverSkcNum   int `json:"deliverSkcNum"`
	DeliveryStatus  int `json:"deliveryStatus"`
}
