package entity

// ShipOrderStaging 发货台
type ShipOrderStaging struct {
	OrderDetailVOList []struct {
		ProductSkuID                int      `json:"productSkuId"`
		ProductSkuImgURLList        []string `json:"productSkuImgUrlList"`
		Color                       string   `json:"color"`
		Size                        string   `json:"size"`
		SkuDeliveryQuantityMaxLimit int      `json:"skuDeliveryQuantityMaxLimit"`
		ProductOriginalSkuID        int      `json:"productOriginalSkuId"`
		ProductSkuPurchaseQuantity  int      `json:"productSkuPurchaseQuantity"`
	} `json:"orderDetailVOList"`
	SubPurchaseOrderBasicVO struct {
		SupplierID                       int    `json:"supplierId"`
		IsCustomProduct                  bool   `json:"isCustomProduct"`
		ExpectLatestArrivalTimeOrDefault any    `json:"expectLatestArrivalTimeOrDefault"`
		ProductSkcPicture                string `json:"productSkcPicture"`
		ProductName                      any    `json:"productName"`
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
		ExpectLatestDeliverTimeOrDefault any    `json:"expectLatestDeliverTimeOrDefault"`
		ArrivalUpcomingDelayTimeMillis   int    `json:"arrivalUpcomingDelayTimeMillis"`
		ReceiveAddressInfo               struct {
			DistrictCode  int    `json:"districtCode"`
			CityName      string `json:"cityName"`
			DistrictName  string `json:"districtName"`
			ProvinceCode  int    `json:"provinceCode"`
			CityCode      int    `json:"cityCode"`
			DetailAddress string `json:"detailAddress"`
			ProvinceName  string `json:"provinceName"`
		} `json:"receiveAddressInfo"`
		AutoRemoveFromDeliveryPlatformTime int    `json:"autoRemoveFromDeliveryPlatformTime"`
		ArrivalDisplayCountdownMillis      int    `json:"arrivalDisplayCountdownMillis"`
		FragileTag                         bool   `json:"fragileTag"`
		PurchaseQuantity                   int    `json:"purchaseQuantity"`
		SubWarehouseName                   string `json:"subWarehouseName"`
		PurchaseTime                       int    `json:"purchaseTime"`
		SubPurchaseOrderSn                 string `json:"subPurchaseOrderSn"`
	} `json:"subPurchaseOrderBasicVO"`
}
