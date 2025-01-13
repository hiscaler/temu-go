package is

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"slices"
)

// RegionId 区域 ID 验证
func RegionId(regionIds []int) validation.RuleFunc {
	return func(value interface{}) error {
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("无效的区域 %v", value)
		}

		if !slices.Contains(regionIds, v) {
			return fmt.Errorf("无效的区域 %v", v)
		}

		return nil
	}
}
