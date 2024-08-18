package temu

import (
	"github.com/hiscaler/temu-go/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_logisticsCompanyService_Companies(t *testing.T) {
	companies, err := temuClient.Services.LogisticsService.Companies()
	assert.Equal(t, nil, err, "Services.LogisticsService.Companies()")
	for _, company := range companies {
		var com entity.LogisticsCompany
		com, err = temuClient.Services.LogisticsService.Company(company.ShipId)
		assert.Equal(t, nil, err, "Services.LogisticsService.Company(%d)", company.ShipId)
		assert.Equal(t, company, com, "Services.LogisticsService.Company(%d)", company.ShipId)
	}
}
