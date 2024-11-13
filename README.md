Temu SDK for Golang
===================

## 服务说明

|   | 服务               | 说明    |
|---|------------------|-------|
| 1 | PurchaseOrder    | 备货单   |
| 2 | ShipOrderStaging | 发货台   |
| 3 | ShipOrder        | 发货单   |

## 使用

```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/hiscaler/temu-go"
	"github.com/hiscaler/temu-go/config"
	"os"
)

func main() {
	b, err := os.ReadFile("./config/config_test.json")
	if err != nil {
		panic(fmt.Sprintf("Read config error: %s", err.Error()))
	}
	var c config.Config
	err = json.Unmarshal(b, &c)
	if err != nil {
		panic(fmt.Sprintf("Parse config file error: %s", err.Error()))
	}
	temuClient := temu.NewClient(c)
}
```

## 文档地址

[TEMU 开发者指南](https://seller.kuajingmaihuo.com/sop/view/634117628601810731)

1. [全托管系统对接指南 - 备货及V3发货](https://seller.kuajingmaihuo.com/sop/view/889973754324016047#YSg2AE)

## 流程

1. 将备货单数据加入发货台；
2. 将发货台数据生成采购单；
3. 打印货物商标、箱唛进行平台发货操作。

## 说明

备货单、采购单是同一个含义，从平台的角度理解是平台向商家下采购单，从商家的角度理解是平台推送过来备货单。

## 注意事项

1. 加入发货台后是不能立即看到物流商数据的