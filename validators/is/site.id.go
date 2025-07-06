package is

import (
	"slices"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// SiteId 站点 ID 验证
func SiteId(siteIds []int) validation.RuleFunc {
	return func(value interface{}) error {
		var err validation.ErrorObject
		v, ok := value.(int)
		if !ok {
			return err.
				SetCode("InvalidSiteId").
				SetParams(map[string]any{"Value": value}).
				SetMessage("无效的站点 {{.Value}}")
		}

		if !slices.Contains(siteIds, v) {
			return err.
				SetCode("InvalidSiteId").
				SetParams(map[string]any{"Value": v}).
				SetMessage("无效的站点 {{.Value}}")
		}

		return nil
	}
}
