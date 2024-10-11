package entity

import "gopkg.in/guregu/null.v4"

// ShipOrderPackageDetail 发货单包裹详情
type ShipOrderPackageDetail struct {
	ProductSkuId         int         `json:"productSkuId"`         // 定制 skuId
	ProductOriginalSkuId int         `json:"productOriginalSkuId"` // skuId
	PersonalText         null.String `json:"personalText"`         // 定制内容
	SkuNum               int         `json:"skuNum"`               // sku 数量
}

// ShipOrderPackage 发货单包裹数据
type ShipOrderPackage struct {
	PackageSn      string                   `json:"packageSn"`      // 包裹号
	ProductSkcId   int                      `json:"productSkcId"`   // skcId
	SkcNum         int                      `json:"skcNum"`         // skc 数量
	PackageDetails []ShipOrderPackageDetail `json:"packageDetails"` // 包裹明细
}
