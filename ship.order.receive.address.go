package temu

import (
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

// 大仓收货地址
type shipOrderReceiveAddressService service

// All [WIP] 查询大仓收货地址 V2
// https://seller.kuajingmaihuo.com/sop/view/889973754324016047#chUUk1
func (s shipOrderReceiveAddressService) All(subPurchaseOrderSnList ...string) (items any, err error) {
	if len(subPurchaseOrderSnList) == 0 {
		return
	}

	var result = struct {
		normal.Response
		Result struct {
			SubPurchaseReceiveAddressGroups []struct {
				SubWarehouseId     int                              `json:"subWarehouseId"`     // 子仓 ID
				ReceiveAddressInfo []entity.ShipOrderReceiveAddress `json:"receiveAddressInfo"` // 收货地址信息
			} `json:"subPurchaseReceiveAddressGroups"` // 子采购单收货地址分组信息列表
			SubPurchaseOrderBasicVOList []string `json:"subPurchaseOrderBasicVOList"` // 子采购单号列表
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
	if len(result.Result.SubPurchaseReceiveAddressGroups) == 0 {
		return
	}
	return result.Result.SubPurchaseReceiveAddressGroups[0], nil
}

// One [WIP] 查询单个备货单收货地址
func (s shipOrderReceiveAddressService) One(subPurchaseOrderSn string) (item any, err error) {
	_, err = s.All(subPurchaseOrderSn)
	if err != nil {
		return
	}

	return
}
