package is

import (
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"regexp"
)

func ShipOrderNumber() validation.RuleFunc {
	return func(value interface{}) error {
		s, ok := value.(string)
		if !ok {
			return errors.New("无效的发货单号。")
		}
		if s == "" {
			return errors.New("发货单号不能为空。")
		}
		if matched, err := regexp.MatchString("^(?i)fh[0-9]{13}$", s); err != nil || !matched {
			return fmt.Errorf("无效的发货单号：%s。", s)
		}
		return nil
	}
}
