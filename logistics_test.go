package temu

import (
	"github.com/hiscaler/temu-go/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_logisticsCompanyService_Companies(t *testing.T) {
	companies, err := temuClient.Services.Logistics.Companies(ctx)
	assert.Equal(t, nil, err, "Services.Logistics.Companies(ctx)")
	for _, company := range companies {
		var com entity.LogisticsShippingCompany
		com, err = temuClient.Services.Logistics.Company(ctx, company.ShipId)
		assert.Equalf(t, nil, err, "Services.Logistics.Company(ctx, %d)", company.ShipId)
		assert.Equalf(t, company, com, "Services.Logistics.Company(ctx, %d)", company.ShipId)
	}
}
