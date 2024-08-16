package temu

import (
	"encoding/json"
	"fmt"
	"github.com/hiscaler/temu-go/config"
	"os"
	"testing"
)

var temuClient *Client

func TestMain(m *testing.M) {
	b, err := os.ReadFile("./config/config_test.json")
	if err != nil {
		panic(fmt.Sprintf("Read config error: %s", err.Error()))
	}
	var c config.Config
	err = json.Unmarshal(b, &c)
	if err != nil {
		panic(fmt.Sprintf("Parse config file error: %s", err.Error()))
	}

	temuClient = NewClient(c)
	m.Run()
}
