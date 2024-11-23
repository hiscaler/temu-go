package temu

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
	"github.com/hiscaler/temu-go/validators/is"
)

// 发货包裹

type shipOrderPackageService service

type ShipOrderPackageQueryParams struct {
	normal.Parameter
	DeliveryOrderSn string `json:"deliveryOrderSn"` // 发货单号
}

func (m ShipOrderPackageQueryParams) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.DeliveryOrderSn,
			validation.Required.Error("发货单号不能为空"),
			validation.By(is.ShipOrderNumber()),
		),
	)
}

// One 发货包裹查询
// https://seller.kuajingmaihuo.com/sop/view/889973754324016047#eprtWq
func (s shipOrderPackageService) One(ctx context.Context, deliveryOrderSn string) (items []entity.ShipOrderPackage, err error) {
	params := ShipOrderPackageQueryParams{DeliveryOrderSn: deliveryOrderSn}
	if err = params.Validate(); err != nil {
		return
	}

	var result = struct {
		normal.Response
		Result struct {
			PackageInfo []entity.ShipOrderPackage `json:"packageInfo"`
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(params).
		SetResult(&result).
		Post("bg.shiporder.package.get")
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	items = result.Result.PackageInfo
	return
}

// 发货包裹编辑

// ShipOrderPackageUpdateRequestDeliverOrderDetail  发货单详情
type ShipOrderPackageUpdateRequestDeliverOrderDetail struct {
	ProductSkuId  int64 `json:"productSkuId"`  // skuId
	DeliverSkuNum int   `json:"deliverSkuNum"` // 发货 sku 数目
}

// ShipOrderPackageUpdateRequestPackageDetail 包裹明细
type ShipOrderPackageUpdateRequestPackageDetail struct {
	ProductSkuId int64 `json:"productSkuId"` // skuId
	SkuNum       int   `json:"skuNum"`       // 发货 sku 数目
}

// ShipOrderPackageUpdateRequestPackage 包裹信息
type ShipOrderPackageUpdateRequestPackage struct {
	PackageDetailSaveInfos []ShipOrderPackageUpdateRequestPackageDetail `json:"packageDetailSaveInfos"` // 包裹明细
}

type ShipOrderPackageUpdateRequest struct {
	normal.Parameter
	DeliveryOrderSn         string                                            `json:"deliveryOrderSn"`         // 发货单号
	DeliverOrderDetailInfos []ShipOrderPackageUpdateRequestDeliverOrderDetail `json:"deliverOrderDetailInfos"` // 发货单详情列表
	PackageInfos            []ShipOrderPackageUpdateRequestPackage            `json:"packageInfos"`            // 包裹信息列表
}

func (m ShipOrderPackageUpdateRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.DeliveryOrderSn,
			validation.Required.Error("发货单号不能为空"),
			validation.By(is.ShipOrderNumber()),
		),
		validation.Field(&m.DeliverOrderDetailInfos, validation.Required.Error("发货单详情列表不能为空")),
		validation.Field(&m.PackageInfos, validation.Required.Error("包裹信息列表不能为空")),
	)
}

// Update 发货包裹编辑
// https://seller.kuajingmaihuo.com/sop/view/889973754324016047#qSU56c
func (s shipOrderPackageService) Update(ctx context.Context, req ShipOrderPackageUpdateRequest) (ok bool, err error) {
	if err = req.Validate(); err != nil {
		return
	}

	var result = struct {
		normal.Response
		Result any `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(req).
		SetResult(&result).
		Post("bg.shiporder.package.edit")
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	return true, nil
}
