package entity

// DeliveryAddress 卖家发货地址
type DeliveryAddress struct {
	ID            int64  `json:"id"`
	MallId        int64  `json:"mallId"`
	ProvinceCode  int64  `json:"provinceCode"`
	ProvinceName  string `json:"provinceName"`
	CityCode      int64  `json:"cityCode"`
	CityName      string `json:"cityName"`
	DistrictCode  int64  `json:"districtCode"`
	DistrictName  string `json:"districtName"`
	AddressLabel  string `json:"addressLabel"`
	AddressDetail string `json:"addressDetail"`
	IsDefault     bool   `json:"isDefault"`
}
