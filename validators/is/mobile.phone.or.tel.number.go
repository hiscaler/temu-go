package is

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"strings"
)

// MobilePhoneOrTelNumber 判断是否为有效的手机或座机号码
func MobilePhoneOrTelNumber() validation.RuleFunc {
	return func(value interface{}) error {
		var err validation.ErrorObject
		s, ok := value.(string)
		if !ok {
			return err.
				SetCode("InvalidMobilePhoneOrTelNumber").
				SetParams(map[string]any{"Number": value}).
				SetMessage("无效的手机/座机号码 {{.Number}}")
		}

		if strings.TrimSpace(s) == "" {
			return err.
				SetCode("MobilePhoneOrTelNumberIsEmpty").
				SetMessage("手机/座机号码不能为空")
		}

		if !mobilePhoneNumberPattern.MatchString(s) && !telNumberPattern.MatchString(s) {
			return err.
				SetCode("InvalidMobilePhoneOrTelNumber").
				SetParams(map[string]any{"Number": value}).
				SetMessage("无效的手机/座机号码 {{.Number}}")
		}

		return nil
	}
}
