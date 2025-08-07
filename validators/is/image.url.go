package is

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// ImageUrl 判断是否为有效的图片链接
func ImageUrl() validation.RuleFunc {
	return func(value interface{}) error {
		var err validation.ErrorObject
		s, ok := value.(string)
		if !ok {
			return err.
				SetCode("InvalidImageUrl").
				SetParams(map[string]any{"Value": value}).
				SetMessage("无效的图片链接 {{.Value}}")
		}

		if !imageUrlPattern.MatchString(s) {
			return err.
				SetCode("InvalidImageUrl").
				SetParams(map[string]any{"Value": s}).
				SetMessage("无效的图片链接 {{.Value}}")
		}

		return nil
	}
}
