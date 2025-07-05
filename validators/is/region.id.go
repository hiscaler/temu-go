package is

import (
	"slices"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// RegionId 区域 ID 验证
func RegionId(regionIds []int) validation.RuleFunc {
	return func(value interface{}) error {
		var err validation.ErrorObject
		v, ok := value.(int)
		if !ok {
			return err.
				SetCode("InvalidRegionId").
				SetParams(map[string]any{"Value": value}).
				SetMessage("无效的区域 {{.Value}}")
		}

		if !slices.Contains(regionIds, v) {
			return err.
				SetCode("InvalidRegionId").
				SetParams(map[string]any{"Value": v}).
				SetMessage("无效的区域 {{.Value}}")
		}

		return nil
	}
}
