package entity

import (
	"encoding/json"
	"testing"
)

func TestProductSkuSupplierPriceDecode(t *testing.T) {
	var price ProductSkuSupplierPrice
	err := json.Unmarshal([]byte(`{"productSkuId":123,"currencyType":"USD","siteSupplierPrices":[{"siteId":100,"supplierPrice":900,"priceReviewStatus":1}]}`), &price)
	if err != nil {
		t.Fatal(err)
	}
	if price.CurrencyType != "USD" || len(price.SiteSupplierPrices) != 1 || price.SiteSupplierPrices[0].PriceReviewStatus.Int64 != 1 {
		t.Fatalf("unexpected price: %+v", price)
	}
}
