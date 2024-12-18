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
		if !ok || len(strings.TrimSpace(s)) < 7 {
			return errors.New("无效的座机号码")
		}

		if !telNumberPattern.MatchString(s) {
			return fmt.Errorf("无效的座机号码：%s", s)
		}

		return nil
	}
}
