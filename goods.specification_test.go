package temu

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_goodsSpecificationService_Create(t *testing.T) {
	type args struct {
		ctx     context.Context
		request GoodsSpecificationCreateRequest
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "color",
			args: args{
				ctx: ctx,
				request: GoodsSpecificationCreateRequest{
					ParentSpecId: 1001,
					SpecName:     "红色",
				}},
			want: 2,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := temuClient.Services.Goods.Specification.Create(tt.args.ctx, tt.args.request)
			if !tt.wantErr(t, err, fmt.Sprintf("Create(%v, %v)", tt.args.ctx, tt.args.request)) {
				return
			}
			assert.Equalf(t, tt.want, got, "Create(%v, %v)", tt.args.ctx, tt.args.request)
		})
	}
}
