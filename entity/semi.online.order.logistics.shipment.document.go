package entity

import (
	"gopkg.in/guregu/null.v4"
)

// SemiOnlineOrderLogisticsShipmentDocument 打印面单
type SemiOnlineOrderLogisticsShipmentDocument struct {
	PackageSn string `json:"packageSn"` // 包裹号
	// 使用如下方式访问
	// 访问地址：这里返回的 url 10 分钟后失效，超时请重新调用 bg.logistics.shipment.document.get 获得新的 url
	// 请求方式：GET
	// 请求 header 包含以下 5 个参数：
	// {
	//	"toa-app-key": ${app_key},
	//	"toa-access-token": ${access_token},
	//	"toa-random": "${32 random numbers}",
	//	"toa-timestamp":"${timestamp}",
	//	"toa-sign":"${sign}"
	// }
	// 其中 toa-random 自定义 32 个随机数
	// 其中 toa-timestamp 长度要求 10 位，精确到秒级，current_time-300<=toa-timestamp<=current_time+300
	// 其中 toa-sign 的计算逻辑与正常请求里公共参数 sign 的计算逻辑一致：
	// - 将 toa-access-token 和 toa-app-key 和 toa-random 和 toa-timestamp 4 个参数进行首字母以 ASCII 方式升序排列 ascii asc，对于相同字母则使用下个字母做二次排序，字母序为从左到右，以此类推
	// - 排序后的结果按照参数名 $key 参数值 $value 的次序进行字符串拼接，拼接处不包含任何字符
	// - 拼接完成的字符串做进一步拼接成 1 个字符串（包含所有kv字符串的长串），并在该长串的头部及尾部分别拼接 app_secret，完成签名字符串的组装
	// - 最后对签名字符串，使用 MD5 算法加密后，得到的 MD5 加密密文后转为大写，即为 toa-sign 值
	// 接口返回文件流。
	Url        string      `json:"url"`         // 返回 url 是加签资源
	ExpireTime int64       `json:"expire_time"` // 过期时间
	Path       null.String `json:"path"`        // 文件下载后本地的保存路径
	Error      null.String `json:"error"`       // 文件下载保存失败的原因
}
