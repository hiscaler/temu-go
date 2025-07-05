package is

import (
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// PackageNumber 包裹号数据验证
func PackageNumber() validation.RuleFunc {
	return func(value interface{}) error {
		var err validation.ErrorObject
		s, ok := value.(string)
		if !ok {
			return err.
				SetCode("InvalidPackageNumber").
				SetParams(map[string]any{"Number": value}).
				SetMessage("无效的包裹号 {{.Number}}")
		}

		if strings.TrimSpace(s) == "" {
			return err.SetCode("PackageNumberIsEmpty").
				SetMessage("包裹号不能为空")
		}

		if !packageNumberPattern.MatchString(s) {
			return err.
				SetCode("InvalidPackageNumber").
				SetParams(map[string]any{"Number": s}).
				SetMessage("无效的包裹号 {{.Number}}")
		}
		return nil
	}
}
