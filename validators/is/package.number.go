package is

import (
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// PackageNumber 包裹号数据验证
func PackageNumber() validation.RuleFunc {
	return func(value interface{}) error {
		s, ok := value.(string)
		if !ok {
			return errors.New("无效的包裹号")
		}
		if s == "" {
			return errors.New("包裹号不能为空")
		}

		if !packageNumberPattern.MatchString(s) {
			return fmt.Errorf("无效的包裹号：%s", s)
		}
		return nil
	}
}
