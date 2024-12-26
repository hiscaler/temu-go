package is

import (
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"strconv"
)

// Millisecond 判断是否为有效的毫秒值
func Millisecond() validation.RuleFunc {
	return func(value interface{}) error {
		s, ok := value.(int64)
		if !ok || s <= 0 {
			return errors.New("无效的毫秒值")
		}

		if !millisecondPattern.MatchString(strconv.Itoa(int(s))) {
			return fmt.Errorf("无效的毫秒值 %d", s)
		}

		return nil
	}
}
