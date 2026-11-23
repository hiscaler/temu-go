package temu

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"
)

func TestShipOrderPackageService_One(t *testing.T) {
	//{
	//	"pubEsContext": "Shipping Label Created",
	//	"operationMove": "100",
	//	"groupTimeZone": "UTC-08:00",
	//	"isDaylightTime": false,
	//	"orderNo": "GFUS01020118918144",
	//	"enContext": "Shipping Label Created",
	//	"thirdWaybillNo": "TEST-GOFO-LIFECYCLE-001",
	//	"id": "544ffaefa3fa47b1b170d8b8caa26c5e",
	//	"customerEventCode": "10012",
	//	"operationTime": "2025-11-20 00:48:11"
	//}
	groupTimeZone := "UTC-08:00"
	operationTime := "2025-11-20 00:48:11"

	// 1. 定义操作时间的布局（不包含时区信息）
	// 注意：这里使用的是基于 2006-01-02 15:04:05 的参考布局
	timeLayout := "2006-01-02 15:04:05"

	// 2. 解析时间字符串（此时解析出的时间默认为 time.UTC 或 Location: Local，取决于具体 Go 版本和环境，但不含实际偏移量信息）
	// 为确保一致性，通常先假定它是没有位置信息的。
	tt, err := time.Parse(timeLayout, operationTime)
	if err != nil {
		log.Fatalf("解析 operationTime 错误: %v", err)
	}

	// 3. 将 groupTimeZone 字符串转换为 time.Location 对象

	// 移除 "UTC-" 或 "UTC+" 前缀，只保留偏移量 "-08:00" 或 "+08:00"
	offsetStr := strings.TrimPrefix(groupTimeZone, "UTC-")
	// 进一步移除可能存在的冒号（time.FixedZone 需要整数秒偏移量）
	// " -08:00" -> "-0800"
	offsetStr = strings.Replace(offsetStr, ":", "", 1)

	// 将 "HHMM" 格式的偏移量转换为秒数
	var hours, minutes int
	_, err = fmt.Sscanf(offsetStr, "%02d%02d", &hours, &minutes)
	if err != nil {
		log.Fatalf("解析时区偏移量错误: %v", err)
	}
	// 计算总偏移秒数 (例如 -8 小时 = -28800 秒)
	offsetSeconds := (hours * 3600) + (minutes * 60)

	// 使用 time.FixedZone 创建一个固定的时区 Location
	// 参数1是时区名称（任意字符串，通常用偏移量表示），参数2是 UTC 偏移秒数
	location := time.FixedZone(groupTimeZone, offsetSeconds)
	// 4. 将解析出的时间点重新关联到正确的时区 Location
	// Time.In() 方法用于将时间点移动到指定的时区，而不是改变底层的时间戳
	tWithZone := time.Date(tt.Year(), tt.Month(), tt.Day(), tt.Hour(), tt.Minute(), tt.Second(), tt.Nanosecond(), location)

	// 或者，如果您确定 time.Parse 默认解析的时间是 Local 时间，
	// 另一种方法是使用 t.In(location) 来查看在那个时区下的表示。
	// 但使用 time.Date 构造新时间点通常更清晰。

	fmt.Printf("原始时间字符串: %s\n", operationTime)
	fmt.Printf("指定时区: %s\n", groupTimeZone)
	fmt.Printf("解析后的带时区时间 (tWithZone): %v\n", tWithZone)
	fmt.Printf("时间戳 (Unix): %d\n", tWithZone.Unix())

	// 转换为 UTC 查看实际绝对时间
	fmt.Printf("转换为 UTC 时间: %v\n", tWithZone.UTC())
	//deliveryOrderSn := "FH2512016181675"
	//packages, err := temuClient.Services.ShipOrder.Package.One(ctx, deliveryOrderSn)
	//_ = packages
	//assert.Nilf(t, err, "temuClient.Services.ShipOrder.Package.One(ctx, %s)", deliveryOrderSn)

}
