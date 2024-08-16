package temu

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

// 发货包裹

type shipOrderPackageService service

type ShipOrderPackageQueryParams struct {
	normal.Parameter
	DeliveryOrderSn string `json:"deliveryOrderSn"` // 发货单号
}

func (m ShipOrderPackageQueryParams) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.DeliveryOrderSn, validation.Required.Error("发货单号不能为空。")),
	)
}

// One List all staging orders
// 发货包裹查询（bg.shiporder.package.get）
func (s shipOrderPackageService) One(deliveryOrderSn string) (data entity.ShipOrderPackage, err error) {
	params := ShipOrderPackageQueryParams{DeliveryOrderSn: deliveryOrderSn}
	if err = params.Validate(); err != nil {
		return
	}

	var result = struct {
		normal.Response
		Result entity.ShipOrderPackage `json:"result"`
	}{}

	resp, err := s.httpClient.R().SetBody(params).SetResult(&result).Post("bg.shiporder.package.get")
	if err == nil {
		err = parseResponse(resp, result.Response)
	}
	if err != nil {
		return
	}

	data = result.Result
	return
}

type ShipOrderPackageEditRequest struct {
	normal.Parameter
	DeliveryOrderSn         string `json:"deliveryOrderSn"` // 发货单号
	DeliverOrderDetailInfos struct {
		DeliverSkuNum int `json:"deliverSkuNum"` // 发货sku数目
		ProductSkuId  int `json:"productSkuId"`  // skuId
	} `json:"deliverOrderDetailInfos"` // 发货单详情列表
	PackageInfos struct {
		PackageDetailSaveInfos struct {
			SkuNum       int `json:"skuNum"`       // 发货sku数目
			ProductSkuId int `json:"productSkuId"` // skuId
		} `json:"packageDetailSaveInfos"` // 包裹明细
	} `json:"packageInfos"` // 包裹信息列表
}

func (m ShipOrderPackageEditRequest) Validate() error {
	return nil
	// return validation.ValidateStruct(&m,
	// 	validation.Field(&m.Request, validation.When(m.Request != nil, validation.By(func(value interface{}) error {
	//
	// 		return nil
	// 	}))),
	// )
}

// Save 发货包裹编辑
// https://seller.kuajingmaihuo.com/sop/view/889973754324016047#qSU56c
func (s shipOrderPackageService) Save(req ShipOrderPackageEditRequest) (ok bool, err error) {
	if err = req.Validate(); err != nil {
		return
	}

	var result = struct {
		normal.Response
		Result any `json:"result"`
	}{}
	resp, err := s.httpClient.R().SetBody(req).SetResult(&result).Post("bg.shiporder.package.edit")
	if err == nil {
		err = parseResponse(resp, result.Response)
	}
	if err != nil {
		return
	}

	return ok, nil
}
