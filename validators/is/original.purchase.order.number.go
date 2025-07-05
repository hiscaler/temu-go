package is

import (
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// OriginalPurchaseOrderNumber 母备货单号数据验证
func OriginalPurchaseOrderNumber() validation.RuleFunc {
	return func(value interface{}) error {
		var err validation.ErrorObject
		s, ok := value.(string)
		if !ok {
			return err.
				SetCode("InvalidOriginalPurchaseOrderNumber").
				SetParams(map[string]any{"Number": value}).
				SetMessage("无效的母备货单号 {{.Number}}")
		}

		if strings.TrimSpace(s) == "" {
			return err.SetCode("OriginalPurchaseOrderNumberIsEmpty").
				SetMessage("母备货单号不能为空")
		}

		if !originalPurchaseOrderNumberPattern.MatchString(s) {
			return err.
				SetCode("InvalidOriginalPurchaseOrderNumber").
				SetParams(map[string]any{"Number": s}).
				SetMessage("无效的母备货单号 {{.Number}}")
		}
		return nil
	}
}
