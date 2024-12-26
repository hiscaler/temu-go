package is

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"time"
)

// TimeRange 判断日期范围是否有有效
// 仅支持 time.DateTime, time.DateOnly, time.TimeOnly 三种格式
func TimeRange(startTime, endTime, timeLayout any) validation.RuleFunc {
	return func(value any) error {
		start, ok := startTime.(string)
		if !ok || start == "" {
			return fmt.Errorf("无效的开始时间 %v", startTime)
		}

		end, ok := endTime.(string)
		if !ok || end == "" {
			return fmt.Errorf("无效的结束时间 %v", endTime)
		}

		layout, ok := timeLayout.(string)
		if !ok {
			return fmt.Errorf("无效的时间格式 %v", timeLayout)
		}

		if layout != time.DateTime && layout != time.DateOnly && layout != time.TimeOnly {
			return fmt.Errorf("无效的时间格式 %s", layout)
		}

		friendlyLayout := ""
		switch layout {
		case time.DateTime:
			friendlyLayout = "YYYY-MM-DD HH:MM:SS"
		case time.DateOnly:
			friendlyLayout = "YYYY-MM-DD"
		default:
			friendlyLayout = "HH:MM:SS"
		}

		err := validation.Validate(start, validation.Date(layout).Error(fmt.Sprintf("无效的开始时间（%s）格式，有效格式为 %s", start, friendlyLayout)))
		if err != nil {
			return err
		}

		err = validation.Validate(end, validation.Date(layout).Error(fmt.Sprintf("无效的结束时间（%s）格式，有效格式为 %s", end, friendlyLayout)))
		if err != nil {
			return err
		}

		sTime, err := time.ParseInLocation(layout, start, time.Local)
		if err != nil {
			return err
		}

		eTime, err := time.ParseInLocation(layout, end, time.Local)
		if err != nil {
			return err
		}

		if eTime.Before(sTime) {
			return fmt.Errorf("结束时间 %s 不能大于开始时间 %s", end, start)
		}

		return nil
	}
}
