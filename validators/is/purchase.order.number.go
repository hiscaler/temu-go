package is

import (
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// PurchaseOrderNumber 备货单号数据验证
func PurchaseOrderNumber() validation.RuleFunc {
	return func(value interface{}) error {
		s, ok := value.(string)
		if !ok {
			return errors.New("无效的备货单号")
		}
		if s == "" {
			return errors.New("备货单号不能为空")
		}
		if !purchaseOrderNumberPattern.MatchString(s) {
			return fmt.Errorf("无效的备货单号：%s", s)
		}
		return nil
	}
}
