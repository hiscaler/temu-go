package is

import (
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"strings"
)

// ShipOrderNumber 发货单号数据验证
func ShipOrderNumber() validation.RuleFunc {
	return func(value interface{}) error {
		s, ok := value.(string)
		if !ok {
			return fmt.Errorf("无效的发货单号：%v", value)
		}

		if strings.TrimSpace(s) == "" {
			return errors.New("发货单号为空")
		}

		if !shipOrderNumberPattern.MatchString(s) {
			return fmt.Errorf("无效的发货单号：%s", s)
		}
		return nil
	}
}
