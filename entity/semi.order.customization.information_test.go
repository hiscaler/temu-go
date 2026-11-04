package entity

import (
	"testing"

	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSemiOrderCustomizationInformationUnmarshal(t *testing.T) {
	const payload = `{
		"orderSn":"PO-1",
		"templateId":null,
		"customizedSvgList":[{"compressedFileUrl":"https://example.test/customized.zip"}]
	}`

	var information SemiOrderCustomizationInformation
	require.NoError(t, json.Unmarshal([]byte(payload), &information))
	assert.Equal(t, "PO-1", information.OrderSn)
	assert.False(t, information.TemplateId.Valid)
	require.Len(t, information.CustomizedSvgList, 1)
	assert.Equal(t, SignatureUrl("https://example.test/customized.zip"), information.CustomizedSvgList[0].CompressedFileUrl)
}
