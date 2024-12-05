package is

import (
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// ShipOrderNumber 发货单号数据验证
func ShipOrderNumber() validation.RuleFunc {
	return func(value interface{}) error {
		s, ok := value.(string)
		if !ok {
			return errors.New("无效的发货单号")
		}
		if s == "" {
			return errors.New("发货单号不能为空")
		}
		if !shipOrderNumberPattern.MatchString(s) {
			return fmt.Errorf("无效的发货单号：%s", s)
		}
		return nil
	}
}
