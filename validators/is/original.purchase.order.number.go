package is

import (
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"strings"
)

// OriginalPurchaseOrderNumber 母备货单号数据验证
func OriginalPurchaseOrderNumber() validation.RuleFunc {
	return func(value interface{}) error {
		s, ok := value.(string)
		if !ok {
			return fmt.Errorf("无效的母备货单号 %v", value)
		}

		if strings.TrimSpace(s) == "" {
			return errors.New("母备货单号为空")
		}

		if !originalPurchaseOrderNumberPattern.MatchString(s) {
			return fmt.Errorf("无效的母备货单号 %s", s)
		}
		return nil
	}
}
