package entity

type ShipOrderPackage struct {
	PackageInfo []struct {
		PackageSn      string `json:"packageSn"`    // 包裹号
		ProductSkcId   int    `json:"productSkcId"` // skcId
		SkcNum         int    `json:"skcNum"`       // skc数量
		PackageDetails []struct {
			ProductSkuId         int    `json:"productSkuId"`         // skuId
			ProductOriginalSkuId int    `json:"productOriginalSkuId"` // 原skuId
			PersonalText         string `json:"personalText"`         // 定制内容
			SkuNum               int    `json:"skuNum"`               // sku数量
		} `json:"packageDetails"` // 包裹明细
	} `json:"packageInfo"` // 创建生成的发货批次号
}
