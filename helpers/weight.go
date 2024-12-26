package helpers

// TruncateWeightValue 换转为整千克数
// Example:
//
//	1g    = 1000g
//	999g  = 1000g
//	1000g = 1000g
//	1001g = 2000g
func TruncateWeightValue(value int64) int64 {
	if value <= 0 {
		return value
	}

	diffValue := value % 1000
	if diffValue == 0 {
		return value
	}
	return value - diffValue + 1000
}
