package entity

// MallAddress 卖家发货地址
type MallAddress struct {
	ID           int    `json:"id"`
	IsDefault    bool   `json:"isDefault"`
	DistrictCode int    `json:"districtCode"`
	Address      string `json:"addressDetail"`
	CityName     string `json:"cityName"`
	DistrictName string `json:"districtName"`
	MallID       int    `json:"mallId"`
	ProvinceCode int    `json:"provinceCode"`
	CityCode     int    `json:"cityCode"`
	ProvinceName string `json:"provinceName"`
	AddressLabel string `json:"addressLabel"`
}
