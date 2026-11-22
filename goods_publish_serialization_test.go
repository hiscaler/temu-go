package temu

import (
	"encoding/json"
	"testing"
)

func TestGoodsCreateSkuResultDecodesSKCID(t *testing.T) {
	var result GoodsCreateResult
	if err := json.Unmarshal([]byte(`{"productId":1,"productSkuList":[{"productSkcId":2,"productSkuId":3,"extCode":"sku"}]}`), &result); err != nil {
		t.Fatal(err)
	}
	if len(result.ProductSkuList) != 1 || result.ProductSkuList[0].ProductSkcId != 2 {
		t.Fatalf("unexpected SKU result: %+v", result.ProductSkuList)
	}
}
