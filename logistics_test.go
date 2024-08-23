package temu

import (
	"github.com/hiscaler/temu-go/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_logisticsCompanyService_Companies(t *testing.T) {
	companies, err := temuClient.Services.Logistics.Companies()
	assert.Equal(t, nil, err, "Services.Logistics.Companies()")
	for _, company := range companies {
		var com entity.LogisticsCompany
		com, err = temuClient.Services.Logistics.Company(company.ShipId)
		assert.Equalf(t, nil, err, "Services.Logistics.Company(%d)", company.ShipId)
		assert.Equalf(t, company, com, "Services.Logistics.Company(%d)", company.ShipId)
	}
}
