package is

import (
	"errors"
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// TimeRange 判断日期范围是否有有效
// 仅支持 time.DateTime, time.DateOnly, time.TimeOnly 三种格式
func TimeRange(startTime, endTime, timeLayout any) validation.RuleFunc {
	return func(value any) error {
		start, ok := startTime.(string)
		if !ok {
			return fmt.Errorf("无效的开始时间 %v", startTime)
		}
		if start == "" {
			return errors.New("开始时间不能为空")
		}

		end, ok := endTime.(string)
		if !ok {
			return fmt.Errorf("无效的结束时间 %v", endTime)
		}
		if end == "" {
			return errors.New("结束时间不能为空")
		}

		layout, ok := timeLayout.(string)
		if !ok {
			return fmt.Errorf("无效的时间格式 %v", timeLayout)
		}

		if layout == "" {
			return errors.New("时间格式不能为空")
		}

		if layout != time.DateTime && layout != time.DateOnly && layout != time.TimeOnly {
			return fmt.Errorf("无效的时间格式 %s", layout)
		}

		friendlyLayout := ""
		example := ""
		switch layout {
		case time.DateTime:
			friendlyLayout = "YYYY-MM-DD HH:MM:SS"
			example = "2000-12-01 08:00:01"
		case time.DateOnly:
			friendlyLayout = "YYYY-MM-DD"
			example = "2000-12-01"
		default:
			friendlyLayout = "HH:MM:SS"
			example = "08:00:01"
		}

		err := validation.Validate(start, validation.Date(layout).Error(fmt.Sprintf("无效的开始时间 %s，有效格式为 %s，例如 %s", start, friendlyLayout, example)))
		if err != nil {
			return err
		}

		err = validation.Validate(end, validation.Date(layout).Error(fmt.Sprintf("无效的结束时间 %s，有效格式为 %s，例如 %s", end, friendlyLayout, example)))
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
			return fmt.Errorf("结束时间 %s 不能小于开始时间 %s", end, start)
		}

		return nil
	}
}
