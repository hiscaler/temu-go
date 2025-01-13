package is

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"slices"
)

// RegionIds 区域 ID 列表验证
func RegionIds(regionIds []int) validation.RuleFunc {
	return func(value interface{}) error {
		ids, ok := value.([]int)
		if !ok {
			return fmt.Errorf("无效的区域 %v", value)
		}

		invalidIds := make([]int, 0)
		for _, id := range ids {
			if !slices.Contains(regionIds, id) && !slices.Contains(invalidIds, id) {
				invalidIds = append(invalidIds, id)
			}
		}
		if len(invalidIds) != 0 {
			return fmt.Errorf("无效的区域 %v", invalidIds)
		}

		return nil
	}
}
