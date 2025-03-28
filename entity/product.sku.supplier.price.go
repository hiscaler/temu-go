package entity

import "gopkg.in/guregu/null.v4"

// ProductSkuSupplierPrice 	货品 sku 供货价
type ProductSkuSupplierPrice struct {
	ProductId          int64   `json:"productId"`     // 货品 ID
	ProductSkcId       int64   `json:"productSkcId"`  //货品skc ID
	ProductSkuId       int64   `json:"productSkuId"`  // 货品sku ID
	SupplierPrice      float64 `json:"supplierPrice"` // 	供货价
	CurrencyType       float64 `json:"currencyType"`  // 币种
	SiteSupplierPrices []struct {
		SiteId            int64    `json:"siteId"`            // 站点 ID
		SupplierPrice     float64  `json:"supplierPrice"`     // 	供货价
		PriceReviewStatus null.Int `json:"priceReviewStatus"` // 核价状态，存量品，或者灰度外可能为空
	} `json:"siteSupplierPrices"` // 站点供货价列表，仅半托管有值
}
