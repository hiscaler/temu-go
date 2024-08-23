package temu

import (
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
	"strings"
)

// 大仓收货地址
type shipOrderReceiveAddressService service

// All 查询大仓收货地址 V2
// https://seller.kuajingmaihuo.com/sop/view/889973754324016047#chUUk1
func (s shipOrderReceiveAddressService) All(subPurchaseOrderSnList ...string) (items []entity.ShipOrderReceiveAddress, err error) {
	if len(subPurchaseOrderSnList) == 0 {
		return
	}

	var result = struct {
		normal.Response
		Result struct {
			SubPurchaseReceiveAddressGroups []entity.ShipOrderReceiveAddress `json:"subPurchaseReceiveAddressGroups"` // 子采购单收货地址分组信息列表
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetBody(map[string][]string{"subPurchaseOrderSnList": subPurchaseOrderSnList}).
		SetResult(&result).
		Post("bg.shiporder.receiveaddressv2.get")
	if err == nil {
		err = parseResponse(resp, result.Response)
	}
	if err != nil {
		return
	}

	return result.Result.SubPurchaseReceiveAddressGroups, nil
}

// One [WIP] 查询单个备货单收货地址
func (s shipOrderReceiveAddressService) One(subPurchaseOrderSn string) (item entity.ShipOrderReceiveAddress, err error) {
	items, err := s.All(subPurchaseOrderSn)
	if err != nil {
		return
	}

	for _, d := range items {
		for _, sn := range d.SubPurchaseOrderSnList {
			if strings.EqualFold(sn, subPurchaseOrderSn) {
				return d, nil
			}
		}
	}
	return item, ErrNotFound
}
