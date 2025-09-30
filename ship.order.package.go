package temu

import (
	"context"
	"errors"

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

func (m ShipOrderPackageQueryParams) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.DeliveryOrderSn,
			validation.Required.Error("发货单号不能为空"),
			validation.By(is.ShipOrderNumber()),
		),
	)
}

// One 发货包裹查询
// https://seller.kuajingmaihuo.com/sop/view/889973754324016047#eprtWq
func (s shipOrderPackageService) One(ctx context.Context, deliveryOrderNumber string) ([]entity.ShipOrderPackage, error) {
	params := ShipOrderPackageQueryParams{DeliveryOrderSn: deliveryOrderNumber}
	if err := params.validate(); err != nil {
		return nil, invalidInput(err)
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
		return nil, err
	}

	return result.Result.PackageInfo, nil
}

// 发货包裹编辑

// ShipOrderPackageUpdateRequestDeliverOrderDetail  发货单详情
type ShipOrderPackageUpdateRequestDeliverOrderDetail struct {
	ProductSkuId  int64 `json:"productSkuId"`  // skuId
	DeliverSkuNum int   `json:"deliverSkuNum"` // 发货 sku 数目
}

func (m ShipOrderPackageUpdateRequestDeliverOrderDetail) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.ProductSkuId, validation.Required.Error("SKU 不能为空")),
		validation.Field(&m.DeliverSkuNum, validation.Min(1).Error("发货数量不能小于 {.min}")),
	)
}

// ShipOrderPackageUpdateRequestPackageDetail 包裹明细
type ShipOrderPackageUpdateRequestPackageDetail struct {
	ProductSkuId int64 `json:"productSkuId"` // skuId
	SkuNum       int   `json:"skuNum"`       // 发货 sku 数目
}

func (m ShipOrderPackageUpdateRequestPackageDetail) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.ProductSkuId, validation.Required.Error("SKU 不能为空")),
		validation.Field(&m.SkuNum, validation.Min(1).Error("发货数量不能小于 {.min}")),
	)
}

// ShipOrderPackageUpdateRequestPackage 包裹信息
type ShipOrderPackageUpdateRequestPackage struct {
	PackageDetailSaveInfos []ShipOrderPackageUpdateRequestPackageDetail `json:"packageDetailSaveInfos"` // 包裹明细
}

func (m ShipOrderPackageUpdateRequestPackage) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.PackageDetailSaveInfos,
			validation.Required.Error("发货包裹不能为空"),
			validation.Each(validation.By(func(value interface{}) error {
				v, ok := value.(ShipOrderPackageUpdateRequestPackageDetail)
				if !ok {
					return errors.New("无效的发货包裹详情")
				}

				return v.validate()
			})),
		),
	)
}

type ShipOrderPackageUpdateRequest struct {
	normal.Parameter
	DeliveryOrderSn         string                                            `json:"deliveryOrderSn"`         // 发货单号
	DeliverOrderDetailInfos []ShipOrderPackageUpdateRequestDeliverOrderDetail `json:"deliverOrderDetailInfos"` // 发货单详情列表
	PackageInfos            []ShipOrderPackageUpdateRequestPackage            `json:"packageInfos"`            // 包裹信息列表
}

func (m ShipOrderPackageUpdateRequest) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.DeliveryOrderSn,
			validation.Required.Error("发货单号不能为空"),
			validation.By(is.ShipOrderNumber()),
		),
		validation.Field(&m.DeliverOrderDetailInfos,
			validation.Required.Error("发货单详情列表不能为空"),
			validation.Each(validation.By(func(value interface{}) error {
				v, ok := value.(ShipOrderPackageUpdateRequestDeliverOrderDetail)
				if !ok {
					return errors.New("无效的发货单详情")
				}

				return v.validate()
			})),
		),
		validation.Field(&m.PackageInfos,
			validation.Required.Error("包裹信息列表不能为空"),
			validation.Each(validation.By(func(value interface{}) error {
				v, ok := value.(ShipOrderPackageUpdateRequestPackage)
				if !ok {
					return errors.New("无效的发货包裹")
				}

				return v.validate()
			})),
		),
	)
}

// Update 发货包裹编辑
// https://seller.kuajingmaihuo.com/sop/view/889973754324016047#qSU56c
func (s shipOrderPackageService) Update(ctx context.Context, req ShipOrderPackageUpdateRequest) (bool, error) {
	if err := req.validate(); err != nil {
		return false, invalidInput(err)
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
		return false, err
	}

	return true, nil
}
