package entity

// ReceiveAddress 地址信息
type ReceiveAddress struct {
	ProvinceCode  int64  `json:"provinceCode"`           // 省份编码
	ProvinceName  string `json:"provinceName"`           // 省
	CityCode      int64  `json:"cityCode"`               // 市编码
	CityName      string `json:"cityName"`               // 市
	DistrictCode  int64  `json:"districtCode"`           // 区编码
	DistrictName  string `json:"districtName"`           // 区
	DetailAddress string `json:"detailAddress"`          // 详细地址
	ReceiverName  string `json:"receiverName,omitempty"` // 收货人
	Phone         string `json:"phone,omitempty"`        // 联系电话
}
