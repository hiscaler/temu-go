package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSensitiveWordChecker_Execute(t *testing.T) {
	tests := []struct {
		name  string
		text  string
		words []string
		want  bool
		want1 []string
	}{
		{"t1", "test b-", []string{"a", "b"}, false, []string{"b"}},
		{"t2", "hello world!", []string{"lo", "wo"}, true, []string{}},
		{"t3", "hello world!wo", []string{"lo", "wo"}, false, []string{"wo"}},
		{"t4", "helLO world!WO", []string{"lo", "wo"}, false, []string{"wo"}},
		{"t5", "hel-LO world!WO", []string{"lo", "wo"}, false, []string{"lo", "wo"}},
		{"t6", "a-b-c world!WO", []string{"abc", "wo"}, false, []string{"abc", "wo"}},
		{"t7", "Ab-C world!WO", []string{"abc", "wo"}, false, []string{"abc", "wo"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			swc := NewSensitiveWordChecker()
			got, got1, err := swc.SetWords(tt.words).Execute(tt.text)
			assert.Nil(t, err)
			assert.Equalf(t, tt.want, got, "Execute(%v)", tt.text)
			assert.Equalf(t, tt.want1, got1, "Execute(%v)", tt.text)
		})
	}
}
