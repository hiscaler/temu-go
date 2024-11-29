package is

import (
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"regexp"
	"strconv"
)

// Millisecond 判断是否为有效的毫秒值
func Millisecond() validation.RuleFunc {
	return func(value interface{}) error {
		s, ok := value.(int64)
		if !ok || s <= 0 {
			return errors.New("无效的毫秒值")
		}

		if matched, err := regexp.MatchString("^[1-9]?[0-9]{12}$", strconv.Itoa(int(s))); err != nil || !matched {
			return fmt.Errorf("无效的毫秒值：%d", s)
		}

		return nil
	}
}
