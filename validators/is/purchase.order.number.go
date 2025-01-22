package is

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"strings"
)

// PurchaseOrderNumber 备货单号数据验证
func PurchaseOrderNumber() validation.RuleFunc {
	return func(value interface{}) error {
		var err validation.ErrorObject
		s, ok := value.(string)
		if !ok {
			return err.
				SetCode("InvalidPurchaseOrder").
				SetParams(map[string]interface{}{"Number": value}).
				SetMessage("无效的备货单号 {{.Number}}")
		}

		if strings.TrimSpace(s) == "" {
			return err.
				SetCode("PurchaseOrderNumberIsEmpty").
				SetMessage("备货单号为空")
		}

		if !purchaseOrderNumberPattern.MatchString(s) {
			return err.
				SetCode("InvalidPurchaseOrder").
				SetParams(map[string]interface{}{"Number": s}).
				SetMessage("无效的备货单号 {{.Number}}")
		}
		return nil
	}

}
