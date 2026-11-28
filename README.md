Temu SDK for Golang
===================

## API 请求地址

| 环境            | url                                              | 使用场景                          |
|-----------------|--------------------------------------------------|-----------------------------------|
| 正式环境 CN     | https://openapi.kuajingmaihuo.com/openapi/router | 全/半托发品，全托寄样、备货和发货 |
| 正式环境 US     | https://openapi-b-us.temu.com/openapi/router     | US半托订单、履约                  |
| 正式环境 EU     | https://openapi-b-eu.temu.com/openapi/router     | EU半托订单、履约                  |
| 正式环境 GLOBAL | https://openapi-b-global.temu.com/openapi/router | GLOBAL半托订单、履约              |
| 正式环境 本本US | https://openapi-b-us.temu.com/openapi/router     | US本本发品、订单和履约全流程      |
| 测试环境 本本US | https://openapi-b-us.temudemo.com/openapi/router | US本本发品、订单和履约全流程      |

## 服务说明

|   | 服务             | 说明   |
|---|------------------|--------|
| 1 | PurchaseOrder    | 备货单 |
| 2 | ShipOrderStaging | 发货台 |
| 3 | ShipOrder        | 发货单 |

## 使用

HTTP 请求默认直连。需要代理时，在 `config.Config` 中设置 `ProxyURL`，或在 JSON 配置中设置
`"proxy_url": "socks5://127.0.0.1:1080"`，也可使用 HTTP 代理地址。该配置与用于替换各区域 API 地址的 `proxies` 相互独立。

将 `log_level` 设置为 `debug` 后，API 请求和响应以 Resty 原生格式通过 `slog` 输出。脱敏规则由调用方在
`config.Config.LogMask` 中指定；未提供规则时不做脱敏，例如 `access_token` 和个人信息会原样输出。

```go
cfg.LogLevel = "debug"
cfg.LogMask = &restyslog.LogMask{
Keep:     1,
Request:  []string{"access_token", "app_secret", "Authorization"},
Response: []string{"phone", "shipping_address.name"},
}
```

示例中的 `restyslog` 对应 `github.com/hiscaler/resty-slog`。`Keep` 表示首尾各保留的字符数，设为 `0` 时每个字符替换为 `*`
，例如 `abc` 变为 `***`；点分路径从 JSON 根节点匹配。规则不会修改实际 HTTP 请求和响应。配置规则后，非 JSON body 和包含原始请求数据的
cURL 会整体隐藏；发送日志前应检查调用方配置的规则是否覆盖所需字段。

`LogLevel` 为空时保留原有 `Debug` 的兼容行为。传入的 `slog.Handler` 也需要允许目标日志级别。`NewClient`
不返回错误，若日志配置无效，会记录配置错误并关闭 HTTP debug 日志。

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

| 服务地址                                  | 说明                 | 查询参数参考地址 |
|-------------------------------------------|----------------------|------------------|
| client.Services.PurchaseOrder.Query       | 备货单查询           |                  |
| client.Services.ShipOrder.Query           | 发货单查询           |                  |
| client.Services.Logistics.Companies       | 物流商               |                  |
| client.Services.Goods.Query               | 商品查询             |                  |
| client.Services.Goods.One                 | 根据商品 SKC ID 查询 |                  |
| client.Services.Goods.Create              | 创建商品             |                  |
| client.Services.Goods.Barcode.NormalGoods | 商品条码查询         |                  |
| client.Services.Goods.Brand.Query         | 查询可绑定的品牌接口 |                  |

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

### 权限说明

## 商品

说明书翻译接口、OpenApi查询半托管商品销售库存、发货单查询、查询说明书翻译结果、货品发布类目属性模板查询、采购单查询、采购备货申请、货品模板查询、货品类目查询、装箱发货校验、尺码组元信息查询、采购单查询v2、物流单号与可用物流公司校验、退货包裹详情查询、文件上传接口、卖家地址查询、创建尺码表、文字转图片-global、子类目查询、查询尺码表模板、bas64图片上传-global、快递公司查询、尺码表可选发布码查询、退供包裹明细列表、查询中文类目映射接口、批量识别牛皮癣图片、说明书语种查询信息、修改物流接口、加入发货台接口、大仓收货地址查询v2、寄样发货、根据尺码表模板创建货品尺码表、内外属性映射、查询发货单v2、销售管理数据查询接口、类目搜索、查询尺码分类接口、查询尺码表元信息、尺码组查询、创建发货单接口v3、发货接口、自行委托三方物流公司查询接口、货品品牌查询、质检结果详情查看、查询尺码表模板、可添加模特类目查询、装箱发货接口、发货包裹查询、寄样单查询、退货包裹查询接口、根据站点查询可绑定的发货仓库信息接口、查询发货台接口、修改备货单下单数量、平台推荐物流商匹配接口、编辑模特信息、质检列表查询、新增模特信息、内外属性映射、箱唛查询、发货单取消、查询脚长对应的鞋码、平台推荐物流商匹配接口V2、发货包裹编辑、卖家发货地址创建、OpenApi半托管销售库存更新接口、查询尺码模板规则、模特信息查询、查询当前token对应授权信息-global、查询当前token对应店铺类型信息、获取预估体积、平台推荐物流商匹配接口V3、物流运单标签获取、批量取消待接单的备货单、定制品商品条码查询、商品条码查询V2、销售管理分仓组数据查询接口、货品车型库模板查询、尺寸图校验、查询货品生命周期状态、发起货品修改单、查询视频转码结果接口-global、货品车型库属性值查询、商品列表查询、上传供应商货品、商品详情查询接口、发起货品修改单、检索OE/OEM数据、卖家中心已废弃列表、best
seller招标单查询、查询视频上传sign接口-global、批量查询爆款售罄商品、货品搬运接口、尺寸图校验结果查询、寄样快递公司查询、半托管新增路由绑定及库存填写接口、货品供货价查询、同意核价单建议价、不同意核价单建议价（并给出新的申报价）、分页查询半托管核价单、分页查询核价单、半托管不同意核价单建议价（并给出新的申报价）、半托管同意核价单建议价、活动报名提交、查询活动报名记录、查询活动列表、查询活动详情、查询活动商品、查询活动场次列表、分页查询全托管调价单、分页查询半托管调价单、分页查询半托管调价单、半托管批量确认/拒绝调价单、半托管批量确认/拒绝调价单、全托管批量确认调价单、商品图片翻译、批量识别牛皮癣图片、图片中cm转inch、高清图片压缩处理、色块图获取、外部商品图片映射temu类目、商品图片翻译接口查询、编辑货品敏感品属性、编辑商品运费模板、新增货品属性、提交货品修改单、修改商品素材、编辑货品属性、货品更新接口、编辑货品尺码表、编辑货品说明书、查询运费模板列表、纠正商品类目、创建规格、查询父规格列表、类目必填信息接口、货品包装清单类型查询、查询商品类目纠正列表、打开JIT、操作日志查询接口、批量修改广告接口、广告投资回报率查询接口、广告商品可创建查询接口、广告投放状态查询接口、创建广告接口、整体投放数据效果（店铺维度）、修改广告接口、广告商品投放数据效果（商品维度）、虚拟库存编辑、虚拟库存查询、编辑属性待办、获取新部件序列号、图片链接转换接口、轮询待办信息、查询包装清单详情页、查询全部包装清单标签、查询待办详情页、查询待办待维护属性、查询包装清单列表、编辑包装清单

## 订单

查询退货物流商信息、上传退货label、处理退款、同意取消订单、查询父售后单信息、查询取消订单售后信息、查询子售后单信息、查询售后单详情、查询售后逆向物流信息、广告投资回报率预估接口、广告投放详情查询接口、广告日志查询接口、广告创建接口、广告整体投放数据效果（店铺维度）接口、广告商品投放数据效果（商品维度）接口、批量修改广告接口、广告修改接口、广告商品可创建查询接口、商品合规标签查询、编辑合规信息接口、查询合规信息元数据接口、为已授权的店铺开通或取消订阅消息服务、获取Access
Token、查询当前token对应授权信息、补货单详情查询接口、创建补货单接口、同步发货单至卖家库存接口、更新补货状态接口、补货计划列表查询接口、预约包裹接口、订单履约信息同步接口、查询包裹预约结果接口、预约单取消接口、查询平台支持发货的物流商列表、可合并发货订单列表、预约包裹接口、订单履约信息同步接口、面单查询接口、查询平台支持发货的物流商列表、temu.track.trackinginfo.get、ERP向合作仓提交履约单接口、ERP查询已建联的合作仓列表接口、ERP未使用合作仓推单原因反馈、ERP履约单信息查询接口、ERP合作仓授权信息同步接口、ERP通知合作仓取消履约单接口