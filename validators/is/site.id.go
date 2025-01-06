package is

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"slices"
)

// SiteId 站点 ID 验证
func SiteId(siteIds []int) validation.RuleFunc {
	return func(value interface{}) error {
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("无效的站点 %v", value)
		}

		if !slices.Contains(siteIds, v) {
			return fmt.Errorf("无效的站点 %v", v)
		}

		return nil
	}
}
