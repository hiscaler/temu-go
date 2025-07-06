package is

import (
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// ShipOrderNumber 发货单号数据验证
func ShipOrderNumber() validation.RuleFunc {
	return func(value interface{}) error {
		var err validation.ErrorObject
		s, ok := value.(string)
		if !ok {
			return err.
				SetCode("InvalidShipOrderNumber").
				SetParams(map[string]any{"Number": value}).
				SetMessage("无效的发货单号 {{.Number}}")
		}

		if strings.TrimSpace(s) == "" {
			return err.SetCode("ShipOrderNumberIsEmpty").
				SetMessage("发货单号不能为空")
		}

		if !shipOrderNumberPattern.MatchString(s) {
			return err.
				SetCode("InvalidShipOrderNumber").
				SetParams(map[string]any{"Number": value}).
				SetMessage("无效的发货单号 {{.Number}}")
		}
		return nil
	}
}
