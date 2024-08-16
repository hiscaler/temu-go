package temu

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_logisticsCompanyService_All(t *testing.T) {
	s, err := temuClient.Services.LogisticsService.Companies()
	fmt.Println(s)
	assert.Equalf(t, nil, err, "test1")
}
