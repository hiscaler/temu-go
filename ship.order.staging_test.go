package temu

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShipOrderStagingService_All(t *testing.T) {
	params := StagingQueryParams{}
	params.Page = 1
	params.PageSize = 10
	items, _, _, _, err := temuClient.Services.ShipOrderStaging.All(params)
	if err != nil {
		t.Errorf("temuClient.Services.ShipOrderStaging.Companies: %s", err.Error())
	} else {
		fmt.Printf("items: %#v", items)
	}
	assert.Equal(t, 4, len(items))
}
