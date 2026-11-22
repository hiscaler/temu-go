Temu SDK for Golang
===================

## API 请求地址

| 环境          | url                                              | 使用场景              |
|-------------|--------------------------------------------------|-------------------|
| 正式环境 CN     | https://openapi.kuajingmaihuo.com/openapi/router | 全/半托发品，全托寄样、备货和发货 |
| 正式环境 US     | https://openapi-b-us.temu.com/openapi/router     | US半托订单、履约         |
| 正式环境 EU     | https://openapi-b-eu.temu.com/openapi/router     | EU半托订单、履约         |
| 正式环境 GLOBAL | https://openapi-b-global.temu.com/openapi/router | GLOBAL半托订单、履约     |
| 正式环境 本本US   | https://openapi-b-us.temu.com/openapi/router     | US本本发品、订单和履约全流程   |
| 测试环境 本本US   | https://openapi-b-us.temudemo.com/openapi/router | US本本发品、订单和履约全流程   |

## 服务说明

|   | 服务               | 说明  |
|---|------------------|-----|
| 1 | PurchaseOrder    | 备货单 |
| 2 | ShipOrderStaging | 发货台 |
| 3 | ShipOrder        | 发货单 |

## 使用

HTTP 请求默认直连。需要代理时，在 `config.Config` 中设置 `ProxyURL`，或在 JSON 配置中设置 `"proxy_url": "socks5://127.0.0.1:1080"`，也可使用 HTTP 代理地址。该配置与用于替换各区域 API 地址的 `proxies` 相互独立。

将 `log_level` 设置为 `debug` 后，API 请求和响应以 Resty 原生格式通过 `slog` 输出。脱敏规则由调用方在 `config.Config.LogMask` 中指定；未提供规则时不做脱敏，例如 `access_token` 和个人信息会原样输出。

```go
cfg.LogLevel = "debug"
cfg.LogMask = &restyslog.LogMask{
    Keep:     1,
    Request:  []string{"access_token", "app_secret", "Authorization"},
    Response: []string{"phone", "shipping_address.name"},
}
```

示例中的 `restyslog` 对应 `github.com/hiscaler/resty-slog`。`Keep` 表示首尾各保留的字符数，设为 `0` 时每个字符替换为 `*`，例如 `abc` 变为 `***`；点分路径从 JSON 根节点匹配。规则不会修改实际 HTTP 请求和响应。配置规则后，非 JSON body 和包含原始请求数据的 cURL 会整体隐藏；发送日志前应检查调用方配置的规则是否覆盖所需字段。

`LogLevel` 为空时保留原有 `Debug` 的兼容行为。传入的 `slog.Handler` 也需要允许目标日志级别。`NewClient` 不返回错误，若日志配置无效，会记录配置错误并关闭 HTTP debug 日志。

本地开发时，可以在 Go workspace 中加入 `resty-slog`。

```go
package main

import (
	"fmt"
	"os"

	"github.com/goccy/go-json"
	"github.com/hiscaler/temu-go"
	"github.com/hiscaler/temu-go/config"
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

## 服务说明

| 服务地址                                      | 说明             | 查询参数参考地址 |
|-------------------------------------------|----------------|----------|
| client.Services.PurchaseOrder.Query       | 备货单查询          |          |
| client.Services.ShipOrder.Query           | 发货单查询          |          |
| client.Services.Logistics.Companies       | 物流商            |          |
| client.Services.Goods.Query               | 商品查询           |          |
| client.Services.Goods.One                 | 根据商品 SKC ID 查询 |          |
| client.Services.Goods.Create              | 创建商品           |          |
| client.Services.Goods.Barcode.NormalGoods | 商品条码查询         |          |
| client.Services.Goods.Brand.Query         | 查询可绑定的品牌接口     |          |

## 参考文档

> [TEMU 开发者指南](https://seller.kuajingmaihuo.com/sop/view/634117628601810731)

1. [全托管系统对接指南 - 备货及V3发货](https://seller.kuajingmaihuo.com/sop/view/889973754324016047#YSg2AE)

> [TEMU OPEN API 更新日志](https://seller.kuajingmaihuo.com/sop/view/512560460535865385)

> [Partner Platform 文档中心](https://partner.kuajingmaihuo.com/document?cataId=875196199516)

> [半托文档](https://partner-us.temu.com/documentation)

> **SMMS 服务商门户**
> 
> [API 文档](https://agentpartner.temu.com/document?cataId=875198836203)
>
> [开发者文档](https://agentpartner.temu.com/document?cataId=875196199516)
> 
> [平台公告](https://agentpartner.temu.com/announcement)

## 全托管备货单

### 流程

1. 将备货单数据加入发货台；
2. 将发货台数据生成采购单；
3. 打印货物商标、箱唛进行平台发货操作。

### 说明

备货单、采购单是同一个含义，从平台的角度理解是平台向商家下采购单，从商家的角度理解是平台推送过来备货单。

### 注意事项

1. 加入发货台后是不能立即看到物流商数据的

## 半托管

### 名词解释

1. PO 单：Parent Order
2. O 单：Order

### 半托官方物流下单发货处理

1. 获取待发货订单
2. 获取发货仓库
3. 根据仓库、订单获取对应的物流服务商
4. 下单
5. 获取面单
6. 第三方仓库执行发货
7. 回写到第三方系统标记发货

## 本本

### 什么是本本

只要主体不是大陆/香港的，都是本本
