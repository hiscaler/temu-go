package is

import (
	"errors"
	"fmt"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// TelNumberAreaCode 判断是否为有效的座机号码区号
func TelNumberAreaCode() validation.RuleFunc {
	return func(value interface{}) error {
		s, ok := value.(string)
		if !ok {
			return fmt.Errorf("无效的座机号码区号 %v", value)
		}

		if strings.TrimSpace(s) == "" {
			return errors.New("座机号码区号为空")
		}

		if !telNumberAreaCodePattern.MatchString(s) {
			return fmt.Errorf("无效的座机号码区号 %s", s)
		}

		return nil
	}
}
