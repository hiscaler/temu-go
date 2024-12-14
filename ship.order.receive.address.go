package temu

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
	"github.com/hiscaler/temu-go/validators/is"
	"strings"
)

// 大仓收货地址
type shipOrderReceiveAddressService service

// Query 查询大仓收货地址 V2
// https://seller.kuajingmaihuo.com/sop/view/889973754324016047#chUUk1
func (s shipOrderReceiveAddressService) Query(ctx context.Context, purchaseOrderNumbers ...string) (items []entity.ShipOrderReceiveAddress, err error) {
	err = validation.Validate(purchaseOrderNumbers,
		validation.Required.Error("备货单号不能为空"),
		validation.Each(validation.By(is.PurchaseOrderNumber())),
	)
	if err != nil {
		return items, invalidInput(err)
	}

	var result = struct {
		normal.Response
		Result struct {
			SubPurchaseReceiveAddressGroups []entity.ShipOrderReceiveAddress `json:"subPurchaseReceiveAddressGroups"` // 子采购单收货地址分组信息列表
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(map[string][]string{"subPurchaseOrderSnList": purchaseOrderNumbers}).
		SetResult(&result).
		Post("bg.shiporder.receiveaddressv2.get")
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	return result.Result.SubPurchaseReceiveAddressGroups, nil
}

// One [WIP] 查询单个备货单收货地址
func (s shipOrderReceiveAddressService) One(ctx context.Context, purchaseOrderNumber string) (item entity.ShipOrderReceiveAddress, err error) {
	items, err := s.Query(ctx, purchaseOrderNumber)
	if err != nil {
		return
	}

	for _, d := range items {
		for _, sn := range d.SubPurchaseOrderSnList {
			if strings.EqualFold(sn, purchaseOrderNumber) {
				return d, nil
			}
		}
	}
	return item, ErrNotFound
}
