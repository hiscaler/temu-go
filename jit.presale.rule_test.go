package temu

import (
	"context"
	"fmt"
	"testing"

	"github.com/hiscaler/temu-go/entity"
	"github.com/stretchr/testify/assert"
)

func Test_jitPresaleRuleService_Query(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	s := temuClient.Services.Jit.PresaleRule
	tests := []struct {
		name     string
		s        jitPresaleRuleService
		args     args
		wantRule entity.JitPresaleRule
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name: "1",
			s:    s,
			args: args{ctx},
			wantRule: entity.JitPresaleRule{
				Version:     3,
				ProtocolUrl: "https://bdl.cdnfe.com/upload_bdl/seller-protocol/86042e8a-1f7f-4e30-b0d3-e325509d0f79.html",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name: "2",
			s:    s,
			args: args{ctx},
			wantRule: entity.JitPresaleRule{
				Version:     0,
				ProtocolUrl: "https://bdl.cdnfe.com/upload_bdl/seller-protocol/86042e8a-1f7f-4e30-b0d3-e325509d0f79.html",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err != nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRule, err := tt.s.Query(tt.args.ctx)
			if !tt.wantErr(t, err, fmt.Sprintf("Query(%v)", tt.args.ctx)) {
				return
			}
			assert.Equalf(t, tt.wantRule, gotRule, "Query(%v)", tt.args.ctx)
		})
	}
}
