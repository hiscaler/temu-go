package is

import (
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"strings"
)

// TelNumberAreaCode 判断是否为有效的座机号码区号
func TelNumberAreaCode() validation.RuleFunc {
	return func(value interface{}) error {
		s, ok := value.(string)
		if !ok {
			return fmt.Errorf("无效的电话号码区号: %v", value)
		}

		if strings.TrimSpace(s) == "" {
			return errors.New("电话号码区号为空")
		}

		if !telNumberAreaCodePattern.MatchString(s) {
			return fmt.Errorf("无效的电话号码区号：%s", s)
		}

		return nil
	}
}
