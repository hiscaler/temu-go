package is

import (
	"errors"
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

		err := validation.Validate(start, validation.Date(layout).Error("无效的开始时间格式"))
		if err != nil {
			return err
		}

		err = validation.Validate(end, validation.Date(layout).Error("无效的结束时间格式"))
		if err != nil {
			return err
		}

		sTime, err := time.Parse(layout, start)
		if err != nil {
			return err
		}

		eTime, err := time.Parse(layout, end)
		if err != nil {
			return err
		}

		if eTime.Before(sTime) {
			return errors.New("结束时间不能大于开始时间")
		}

		return nil
	}
}
