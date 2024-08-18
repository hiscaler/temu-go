package entity

// ShipOrderReceiveAddress 发货单收货地址
type ShipOrderReceiveAddress struct {
	SubWarehouseId     int `json:"subWarehouseId"` // 子仓 ID
	ReceiveAddressInfo struct {
		DistrictCode  int    `json:"districtCode"`  // 区编码
		DistrictName  string `json:"districtName"`  // 区
		CityCode      int    `json:"cityCode"`      // 市编码
		CityName      string `json:"cityName"`      // 市
		ProvinceCode  int    `json:"provinceCode"`  // 省份编码
		ProvinceName  string `json:"provinceName"`  // 省
		Phone         string `json:"phone"`         // 联系电话
		ReceiverName  string `json:"receiverName"`  // 收货人
		DetailAddress string `json:"detailAddress"` // 详细地址
	} `json:"receiveAddressInfo"` // 收货地址信息
	SubPurchaseOrderBasicVOList []any `json:"subPurchaseOrderBasicVOList"` // 子采购单号列表
}
