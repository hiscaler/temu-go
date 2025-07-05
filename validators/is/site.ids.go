package is

import (
	"fmt"
	"slices"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// SiteIds 站点 ID 列表验证
func SiteIds(siteIds []int) validation.RuleFunc {
	return func(value interface{}) error {
		ids, ok := value.([]int)
		if !ok {
			return fmt.Errorf("无效的站点 %v", value)
		}

		invalidIds := make([]int, 0)
		for _, id := range ids {
			if !slices.Contains(siteIds, id) && !slices.Contains(invalidIds, id) {
				invalidIds = append(invalidIds, id)
			}
		}
		if len(invalidIds) != 0 {
			return fmt.Errorf("无效的站点 %v", invalidIds)
		}

		return nil
	}
}
