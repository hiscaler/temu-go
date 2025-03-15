package temu

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_mallService_Type(t *testing.T) {
	_, err := temuClient.Services.Mall.Type(ctx)
	assert.Equal(t, nil, err, "Services.Mall.Type(ctx)")
}

func Test_mallService_Permission(t *testing.T) {
	v, err := temuClient.Services.Mall.Permission(ctx)
	assert.Equal(t, nil, err, "Services.Mall.Permission(ctx)")
	assert.Equal(t, true, v.Valid(), "Services.Mall.Permission(ctx) access token valid")
}
