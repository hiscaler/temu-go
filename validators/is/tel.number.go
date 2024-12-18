package is

import (
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"strings"
)

// TelNumber 判断是否为有效的座机号码
func TelNumber() validation.RuleFunc {
	return func(value interface{}) error {
		s, ok := value.(string)
		if !ok {
			return fmt.Errorf("无效的座机号码：%v", value)
		}

		if strings.TrimSpace(s) == "" {
			return errors.New("座机号码为空")
		}

		if !telNumberPattern.MatchString(s) {
			return fmt.Errorf("无效的座机号码：%s", s)
		}

		return nil
	}
}
