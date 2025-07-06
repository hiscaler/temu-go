package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizer_ProductName(t *testing.T) {
	type fields struct {
		text            string
		trimSpace       bool // 删除两端的空格
		cleanExtraSpace bool // 清理多余空格
		halfWidth       bool // 是否为半角
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"1",
			fields{
				text:            "a",
				trimSpace:       false,
				cleanExtraSpace: false,
				halfWidth:       false,
			},
			"a",
		},
		{
			"2",
			fields{
				text:            "a   b",
				trimSpace:       false,
				cleanExtraSpace: true,
				halfWidth:       false},
			"a b",
		},
		{
			"3",
			fields{
				text:            " this is    an a,    not an   “b”   , do  you know！",
				trimSpace:       true,
				cleanExtraSpace: true,
				halfWidth:       true,
			},
			`this is an a, not an "b", do you know!`,
		},
	}
	normalizer := NewTextNormalizer()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			normalizer.
				SetText(tt.fields.text).
				TrimSpace(tt.fields.trimSpace).
				CleanExtraSpace(tt.fields.cleanExtraSpace).
				HalfWidth(tt.fields.halfWidth)
			assert.Equalf(t, tt.want, normalizer.String(), "String()")
		})
	}
}
