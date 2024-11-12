package entity

// MallAddress 卖家发货地址
type MallAddress struct {
	ID            int64  `json:"id"`
	IsDefault     bool   `json:"isDefault"`
	DistrictCode  int    `json:"districtCode"`
	AddressDetail string `json:"addressDetail"`
	CityName      string `json:"cityName"`
	DistrictName  string `json:"districtName"`
	MallId        int64  `json:"mallId"`
	ProvinceCode  int    `json:"provinceCode"`
	CityCode      int    `json:"cityCode"`
	ProvinceName  string `json:"provinceName"`
	AddressLabel  string `json:"addressLabel"`
}
