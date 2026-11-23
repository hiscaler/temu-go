package temu

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGoodsCreateGoodsLayerDecorationOmitsDefaultFontFamily(t *testing.T) {
	var decoration GoodsCreateGoodsLayerDecoration
	require.NoError(t, json.Unmarshal([]byte(`{
		"lang": "en",
		"type": "text",
		"priority": 1,
		"key": "DecImage",
		"contentList": [{
			"text": "Product description",
			"textModuleDetails": {
				"backgroundColor": "#ffffff",
				"fontSize": 12,
				"align": "left",
				"fontColor": "#333333"
			}
		}]
	}`), &decoration))

	payload, err := json.Marshal(decoration)

	require.NoError(t, err)
	assert.NotContains(t, string(payload), `"fontFamily"`)
	assert.NotContains(t, string(payload), `"imgUrl"`)
	assert.NotContains(t, string(payload), `"width"`)
	assert.NotContains(t, string(payload), `"height"`)
}

func TestGoodsCreateGoodsLayerDecorationOmitsTextFieldsForImage(t *testing.T) {
	decoration := GoodsCreateGoodsLayerDecoration{
		Lang:     "en",
		Type:     "image",
		Priority: 1,
		Key:      "DecImage",
		ContentList: []GoodsCreateGoodsLayerDecorationContent{{
			ImgUrl: "https://cdn.example.com/detail.jpg",
		}},
	}

	payload, err := json.Marshal(decoration)

	require.NoError(t, err)
	assert.NotContains(t, string(payload), `"text"`)
	assert.NotContains(t, string(payload), `"textModuleDetails"`)
	assert.NotContains(t, string(payload), `"width"`)
	assert.NotContains(t, string(payload), `"height"`)
}
