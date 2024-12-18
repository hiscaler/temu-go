package is

import (
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"strings"
)

// PackageNumber 包裹号数据验证
func PackageNumber() validation.RuleFunc {
	return func(value interface{}) error {
		s, ok := value.(string)
		if !ok {
			return fmt.Errorf("无效的包裹号：%v", value)
		}

		if strings.TrimSpace(s) == "" {
			return errors.New("包裹号为空")
		}

		if !packageNumberPattern.MatchString(s) {
			return fmt.Errorf("无效的包裹号：%s", s)
		}
		return nil
	}
}
