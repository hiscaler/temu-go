package entity

// MallAddress 卖家发货地址
type MallAddress struct {
	ID            int64  `json:"id"`
	IsDefault     bool   `json:"isDefault"`
	DistrictCode  int64  `json:"districtCode"`
	AddressDetail string `json:"addressDetail"`
	CityName      string `json:"cityName"`
	DistrictName  string `json:"districtName"`
	MallID        int    `json:"mallId"`
	ProvinceCode  int64  `json:"provinceCode"`
	CityCode      int64  `json:"cityCode"`
	ProvinceName  string `json:"provinceName"`
	AddressLabel  string `json:"addressLabel"`
}
