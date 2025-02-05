package is

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"slices"
)

// RegionIds 区域 ID 列表验证
func RegionIds(regionIds []int) validation.RuleFunc {
	return func(value interface{}) error {
		var err validation.ErrorObject
		ids, ok := value.([]int)
		if !ok {
			return err.
				SetCode("InvalidRegionId").
				SetParams(map[string]any{"Value": value}).
				SetMessage("无效的区域 {{.Value}}")
		}

		invalidIds := make([]int, 0)
		for _, id := range ids {
			if !slices.Contains(regionIds, id) && !slices.Contains(invalidIds, id) {
				invalidIds = append(invalidIds, id)
			}
		}
		if len(invalidIds) != 0 {
			return err.
				SetCode("InvalidRegionId").
				SetParams(map[string]any{"Value": invalidIds}).
				SetMessage("无效的区域 {{.Value}}")
		}

		return nil
	}
}
