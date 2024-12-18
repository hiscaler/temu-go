package temu

import (
	"context"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
	"github.com/hiscaler/temu-go/validators/is"
)

// 卖家发货地址服务

type mallAddressService service

// Query 卖家发货地址查询
// https://seller.kuajingmaihuo.com/sop/view/889973754324016047#1qow2K
func (s mallAddressService) Query(ctx context.Context) (items []entity.MallAddress, err error) {
	var result = struct {
		normal.Response
		Result []entity.MallAddress `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetResult(&result).
		Post("bg.mall.address.get")
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	return result.Result, nil
}

// One 根据 ID 查询单个卖家发货地址
func (s mallAddressService) One(ctx context.Context, id int64) (address entity.MallAddress, err error) {
	items, err := s.Query(ctx)
	if err != nil {
		return
	}

	for _, item := range items {
		if item.ID == id {
			return item, nil
		}
	}

	return address, ErrNotFound
}

// 卖家发货地址创建
// https://seller.kuajingmaihuo.com/sop/view/889973754324016047#gcyXKJ

type CreateDeliveryAddressRequest struct {
	WarehouseType            int64  `json:"warehouseType"`                      // 仓库类型
	WarehouseAreaType        string `json:"warehouseAreaType"`                  // 仓库面积类型
	ProvinceCode             int64  `json:"provinceCode"`                       // 省份编码
	ProvinceName             string `json:"provinceName"`                       // 省名
	CityCode                 int64  `json:"cityCode"`                           // 市编码
	CityName                 string `json:"cityName"`                           // 市名
	DistrictCode             int64  `json:"districtCode"`                       // 区编码
	DistrictName             string `json:"districtName"`                       // 区名
	TownCode                 int64  `json:"townCode,omitempty"`                 // 城镇编码
	TownName                 string `json:"townName,omitempty"`                 // 城镇
	ContactPersonName        string `json:"contactPersonName"`                  // 联系人
	ContactPersonPhoneAreaNo string `json:"contactPersonPhoneAreaNo,omitempty"` // 联系人电话区号
	ContactPersonPhone       string `json:"contactPersonPhone"`                 // 联系人电话
	AddressLabel             string `json:"addressLabel"`                       // 地址标签
	AddressDetail            string `json:"addressDetail"`                      // 详细地址
}

func (m CreateDeliveryAddressRequest) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.WarehouseType, validation.Required.Error("仓库类型不能为空")),
		validation.Field(&m.WarehouseAreaType, validation.Required.Error("仓库面积类型不能为空")),
		validation.Field(&m.ProvinceCode, validation.Required.Error("省份编码不能为空")),
		validation.Field(&m.ProvinceName, validation.Required.Error("省名不能为空")),
		validation.Field(&m.CityCode, validation.Required.Error("市编码不能为空")),
		validation.Field(&m.CityName, validation.Required.Error("市名不能为空")),
		validation.Field(&m.DistrictCode, validation.Required.Error("区编码不能为空")),
		validation.Field(&m.DistrictName, validation.Required.Error("区名不能为空")),
		validation.Field(&m.ContactPersonName, validation.Required.Error("联系人不能为空")),
		validation.Field(&m.ContactPersonPhone,
			validation.Required.Error("联系人电话不能为空"),
			validation.By(func(value interface{}) error {
				s, ok := value.(string)
				if !ok {
					return fmt.Errorf("无效的联系人电话：%v", s)
				}

				var err error
				if err = validation.Validate(s, validation.By(is.MobilePhoneNumber())); err != nil {
					err = validation.Validate(s, validation.By(is.TelNumber()))
				}
				if err != nil {
					return fmt.Errorf("无效的联系人电话：%s", s)
				}

				return nil
			}),
		),
		validation.Field(&m.AddressLabel, validation.Required.Error("地址标签不能为空")),
		validation.Field(&m.AddressDetail, validation.Required.Error("详细地址不能为空")),
	)
}

func (s mallAddressService) Create(ctx context.Context, request CreateDeliveryAddressRequest) (addressId int64, err error) {
	if err = request.validate(); err != nil {
		err = invalidInput(err)
		return
	}
	var result = struct {
		normal.Response
		Result struct {
			AddressId int64 `json:"addressId"` // 创建的地址 ID
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(request).
		SetResult(&result).
		Post("bg.mall.address.add")
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	return result.Result.AddressId, nil
}
