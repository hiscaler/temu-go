package temu

import (
	"context"
	"fmt"
	"testing"

	"github.com/hiscaler/temu-go/entity"
	"github.com/stretchr/testify/assert"
)

func Test_semiOnlineOrderLogisticsShipmentService_Query(t *testing.T) {
	type args struct {
		ctx            context.Context
		packageNumbers []string
	}
	tests := []struct {
		name    string
		args    args
		want    []entity.SemiOnlineOrderLogisticsShipmentPackageResult
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "1",
			args: args{
				ctx:            ctx,
				packageNumbers: []string{"PK-5265463056138373225"},
			},
			want: []entity.SemiOnlineOrderLogisticsShipmentPackageResult{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err, i)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := temuClient.Services.SemiManaged.OnlineOrder.Logistics.Shipment.Query(ctx, tt.args.packageNumbers...)
			if !tt.wantErr(t, err, fmt.Sprintf("Query(%v, %v)", tt.args.ctx, tt.args.packageNumbers)) {
				return
			}
			assert.Equalf(t, tt.want, got, "Query(%v, %v)", tt.args.ctx, tt.args.packageNumbers)
		})
	}
}
