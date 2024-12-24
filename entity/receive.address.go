package entity

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/temu-go/validators/is"
)

// ReceiveAddress 收货地址信息
type ReceiveAddress struct {
	ProvinceCode  int64  `json:"provinceCode"`    // 省份编码
	ProvinceName  string `json:"provinceName"`    // 省份名称
	CityCode      int64  `json:"cityCode"`        // 市编码
	CityName      string `json:"cityName"`        // 市名称
	DistrictCode  int64  `json:"districtCode"`    // 区编码
	DistrictName  string `json:"districtName"`    // 区名称
	DetailAddress string `json:"detailAddress"`   // 详细地址
	ReceiverName  string `json:"receiverName"`    // 收货人
	Phone         string `json:"phone,omitempty"` // 联系电话
}

func (m ReceiveAddress) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.DetailAddress,
			validation.Required.Error("详细地址不能为空"),
		),
		validation.Field(&m.ReceiverName,
			validation.Required.Error("收货人不能为空"),
		),
		validation.Field(&m.Phone,
			validation.Required.Error("联系电话不能为空"),
			validation.By(is.MobilePhoneOrTelNumber()),
		),
	)
}
