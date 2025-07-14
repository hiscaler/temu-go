package helpers

import (
	"strconv"
	"time"
)

// StrTime2UnixMilli 将字符串时间格式的时间转换为毫秒
func StrTime2UnixMilli(startTime, endTime string) (start, end string, err error) {
	t, err := time.ParseInLocation(time.DateTime, startTime, time.Local)
	if err != nil {
		return
	}
	start = strconv.Itoa(int(t.UnixMilli()))
	t, err = time.ParseInLocation(time.DateTime, endTime, time.Local)
	if err != nil {
		return
	}
	return start, strconv.Itoa(int(t.UnixMilli())), nil
}

// StrTime2Unix 将字符串时间格式的时间转换为秒
func StrTime2Unix(startTime, endTime string) (start, end string, err error) {
	t, err := time.ParseInLocation(time.DateTime, startTime, time.Local)
	if err != nil {
		return
	}
	start = strconv.Itoa(int(t.Unix()))
	t, err = time.ParseInLocation(time.DateTime, endTime, time.Local)
	if err != nil {
		return
	}
	return start, strconv.Itoa(int(t.Unix())), nil
}
