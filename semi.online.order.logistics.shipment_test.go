package temu

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_semiOnlineOrderLogisticsShipmentService_Query(t *testing.T) {
	type args struct {
		ctx            context.Context
		packageNumbers []string
	}
	tests := []struct {
		name    string
		s       semiOnlineOrderLogisticsShipmentService
		args    args
		want    []entity.SemiOnlineOrderLogisticsShipmentPackageResult
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Query(tt.args.ctx, tt.args.packageNumbers...)
			if !tt.wantErr(t, err, fmt.Sprintf("Query(%v, %v)", tt.args.ctx, tt.args.packageNumbers...)) {
				return
			}
			assert.Equalf(t, tt.want, got, "Query(%v, %v)", tt.args.ctx, tt.args.packageNumbers...)
		})
	}
}
