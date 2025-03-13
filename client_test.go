package temu

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/hiscaler/temu-go/config"
	"github.com/hiscaler/temu-go/entity"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
	"os"
	"testing"
)

var temuClient *Client
var ctx context.Context

func TestMain(m *testing.M) {
	b, err := os.ReadFile("./config/config.json")
	if err != nil {
		panic(fmt.Sprintf("Read config error: %s", err.Error()))
	}
	var cfg config.Config
	err = json.Unmarshal(b, &cfg)
	if err != nil {
		panic(fmt.Sprintf("Parse config file error: %s", err.Error()))
	}

	temuClient = NewClient(cfg)
	temuClient.SetLanguage(language.Chinese)
	ctx = context.Background()
	m.Run()
}

func TestClient_SetRegion(t *testing.T) {
	tests := []struct {
		name         string
		region       string
		expectRegion string
	}{
		{"t1", entity.ChinaRegion, entity.ChinaRegion},
		{"t2", entity.AmericanRegion, entity.AmericanRegion},
	}
	for _, tt := range tests {
		temuClient.SetRegion(tt.region)
		assert.Equalf(t, tt.expectRegion, temuClient.Region, tt.name)
	}
}

func Test_getVersion(t *testing.T) {
	tests := []struct {
		name string
		typ  string
		want string
	}{
		{"1", "", "V1"},
		{"2", "a.v1", "V1"},
		{"3", ".v1.", "V1"},
		{"4", "a.v111.b", "V111"},
		{"5", "a.va.b", "V1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, getVersion(tt.typ), "getVersion(%v)", tt.typ)
		})
	}
}
