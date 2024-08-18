package temu

import (
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

// 卖家发货地址服务

type mallAddressService service

// All 卖家发货地址查询
// https://seller.kuajingmaihuo.com/sop/view/889973754324016047#1qow2K
func (s mallAddressService) All() (items []entity.MallAddress, err error) {
	var result = struct {
		normal.Response
		Result []entity.MallAddress `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetResult(&result).
		Post("bg.mall.address.get")
	if err == nil {
		err = parseResponse(resp, result.Response)
	}
	if err != nil {
		return
	}

	return result.Result, nil
}

// One 根据 ID 查询单个卖家发货地址
func (s mallAddressService) One(id int) (address entity.MallAddress, err error) {
	items, err := s.All()
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
