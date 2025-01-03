package entity

// SemiOnlineOrderLogisticsShipmentDocument 打印面单
type SemiOnlineOrderLogisticsShipmentDocument struct {
	PackageSn string `json:"packageSn"` // 包裹号
	// 使用如下方式访问
	//	访问地址：这里返回的url，url10分钟后失效，超时请重新调用bg.logistics.shipment.document.get获得新的url
	//	请求方式：GET
	//	请求header包含以下5个参数：
	// {
	//	"toa-app-key": ${app_key},
	//	"toa-access-token": ${access_token},
	//	"toa-random": "${32 random numbers}",
	//	"toa-timestamp":"${timestamp}",
	//	"toa-sign":"${sign}"
	// }
	//	其中toa-random自定义32个随机数
	//	其中toa-timestamp长度要求10位，精确到秒级，current_time-300<=toa-timestamp<=current_time+300
	//	其中toa-sign的计算逻辑与正常请求里公共参数sign的计算逻辑一致：
	//	- 将toa-access-token和toa-app-key和toa-random和toa-timestamp 4个参数进行首字母以ASCII方式升序排列ascii asc ，对于相同字母则使用下个字母做二次排序，字母序为从左到右，以此类推
	//	- 排序后的结果按照参数名$key 参数值$value 的次序进行字符串拼接，拼接处不包含任何字符
	//	- 拼接完成的字符串做进一步拼接成1个字符串（包含所有kv字符串的长串），并在该长串的头部及尾部分别拼接app_secret ，完成签名字符串的组装
	//	- 最后对签名字符串，使用MD5 算法加密后，得到的MD5 加密密文后转为大写，即为toa-sign 值
	//	接口返回文件流。
	Url string `json:"url"` // 返回 url 是加签资源
}
