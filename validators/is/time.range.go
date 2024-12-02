package is

import (
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"time"
)

// TimeRange 判断日期范围是否有有效
func TimeRange(startTime, endTime, timeLayout any) validation.RuleFunc {
	return func(value any) error {
		start, ok := startTime.(string)
		if !ok || start == "" {
			return errors.New("无效的开始时间")
		}

		end, ok := endTime.(string)
		if !ok || end == "" {
			return errors.New("无效的结束时间")
		}

		layout, ok := timeLayout.(string)
		if !ok {
			return errors.New("无效的时间格式")
		}

		err := validation.Validate(start, validation.Date(layout).Error(fmt.Sprintf("无效的开始时间 %s 格式，有效格式为：%s", start, layout)))
		if err != nil {
			return err
		}

		err = validation.Validate(end, validation.Date(layout).Error(fmt.Sprintf("无效的结束时间 %s 格式，有效格式为：%s", end, layout)))
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
			return errors.New("结束时间不能大于开始时间")
		}

		return nil
	}
}
