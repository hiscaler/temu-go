package helpers

import (
	"strconv"
	"time"
)

// StrTime2UnixMilli 将字符串时间格式的时间转换为毫秒
func StrTime2UnixMilli(startTime, endTime string) (start, end string, err error) {
	t, err := time.ParseInLocation(time.DateTime, startTime+" 00:00:00", time.Local)
	if err != nil {
		return
	}
	start = strconv.Itoa(int(t.UnixMilli()))
	t, err = time.ParseInLocation(time.DateTime, endTime+" 23:59:59", time.Local)
	if err != nil {
		return
	}
	return start, strconv.Itoa(int(t.UnixMilli())), nil
}
