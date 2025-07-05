package is

import (
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// MobilePhoneNumber 判断是否为有效的手机号码
func MobilePhoneNumber() validation.RuleFunc {
	return func(value interface{}) error {
		var err validation.ErrorObject
		s, ok := value.(string)
		if !ok {
			return err.
				SetCode("InvalidMobilePhoneNumber").
				SetParams(map[string]any{"Number": value}).
				SetMessage("无效的手机号码 {{.Number}}")
		}

		if strings.TrimSpace(s) == "" {
			return err.
				SetCode("MobilePhoneNumberIsEmpty").
				SetMessage("手机号码不能为空")
		}

		if !mobilePhoneNumberPattern.MatchString(s) {
			return err.
				SetCode("InvalidMobilePhoneNumber").
				SetParams(map[string]any{"Number": s}).
				SetMessage("无效的手机号码 {{.Number}}")
		}

		return nil
	}
}
