package is

import (
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"strings"
)

// MobilePhoneNumber 判断是否为有效的手机号码
func MobilePhoneNumber() validation.RuleFunc {
	return func(value interface{}) error {
		s, ok := value.(string)
		if !ok || len(strings.TrimSpace(s)) < 11 {
			return errors.New("无效的手机号码")
		}

		if !mobilePhoneNumberPattern.MatchString(s) {
			return fmt.Errorf("无效的手机号码：%s", s)
		}

		return nil
	}
}
