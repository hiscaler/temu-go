package is

import (
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Millisecond 判断是否为有效的毫秒值
func Millisecond() validation.RuleFunc {
	return func(value interface{}) error {
		var err validation.ErrorObject
		s, ok := value.(int64)
		if !ok || s <= 0 {
			return err.
				SetCode("InvalidMillisecondValue").
				SetParams(map[string]any{"Value": value}).
				SetMessage("无效的毫秒值 {{.Value}}")
		}

		if !millisecondPattern.MatchString(strconv.Itoa(int(s))) {
			return err.
				SetCode("InvalidMillisecondValue").
				SetParams(map[string]any{"Value": s}).
				SetMessage("无效的毫秒值 {{.Value}}")
		}

		return nil
	}
}
