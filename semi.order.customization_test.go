package temu

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateCustomizationInformationOrderNumbers(t *testing.T) {
	orderNumbers := make([]string, maxCustomizationInformationOrders+1)
	for i := range orderNumbers {
		orderNumbers[i] = "order"
	}

	err := validateCustomizationInformationOrderNumbers(orderNumbers)
	assert.ErrorContains(t, err, "不能超过 10 个")
}

func TestValidateCustomizationInformationOrderNumbersRejectsBlankOrderNumber(t *testing.T) {
	err := validateCustomizationInformationOrderNumbers([]string{"order", " "})
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "无效的子单号"))
}
