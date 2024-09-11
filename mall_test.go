package temu

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_mallService_IsSemiManaged(t *testing.T) {
	_, err := temuClient.Services.Mall.IsSemiManaged(ctx)
	assert.Equal(t, nil, err, "Services.Mall.IsSemiManaged(ctx")
}

func Test_mallService_Permission(t *testing.T) {
	_, err := temuClient.Services.Mall.Permission(ctx)
	assert.Equal(t, nil, err, "Services.Mall.Permission(ctx")
}
