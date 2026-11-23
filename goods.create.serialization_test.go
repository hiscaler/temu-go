package temu

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGoodsCreateProductSkuOmitsFullManagedPriceForSemiManaged 验证半托管发品不会传零值全托供货价。
func TestGoodsCreateProductSkuOmitsFullManagedPriceForSemiManaged(t *testing.T) {
	payload, err := json.Marshal(GoodsCreateProductSku{
		SiteSupplierPrices: []GoodsCreateProductSkuSiteSupplierPrice{{
			SiteId:        100,
			SupplierPrice: 900,
		}},
	})

	require.NoError(t, err)
	var data map[string]any
	require.NoError(t, json.Unmarshal(payload, &data))
	assert.NotContains(t, data, "supplierPrice")
	assert.Equal(t, []any{map[string]any{"siteId": float64(100), "supplierPrice": float64(900)}}, data["siteSupplierPrices"])
}
