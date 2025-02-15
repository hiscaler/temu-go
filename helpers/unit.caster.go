package helpers

import "math"

// UnitCaster 计量单位转换
// 四舍六入五成双是一种比较精确比较科学的计数保留法，是一种数字修约规则，又名银行家舍入法。它比通常用的四舍五入法更加精确。
// 具体规则为：
// 四舍六入五考虑，五后非零就进一，五后为零看奇偶，五前为偶应舍去，五前为奇要进一
type UnitCaster struct {
	BeforeValue float64
	AfterValue  float64
	Precision   int
}

func NewUnitCaster(precision int) *UnitCaster {
	return &UnitCaster{Precision: precision}
}

// Cm2In 厘米转英寸
func (uc *UnitCaster) Cm2In(cm float64) *UnitCaster {
	uc.BeforeValue = cm
	uc.AfterValue = cm * 0.393701
	return uc
}

func (uc *UnitCaster) In2Cm(in float64) *UnitCaster {
	uc.BeforeValue = in
	uc.AfterValue = in / 0.393701
	return uc
}

func (uc *UnitCaster) G2Kg(g float64) *UnitCaster {
	uc.BeforeValue = g
	uc.AfterValue = g / 1000
	return uc
}

func (uc *UnitCaster) Kg2G(kg float64) *UnitCaster {
	uc.BeforeValue = kg
	uc.AfterValue = kg * 1000
	return uc
}

func (uc *UnitCaster) G2Lb(g float64) *UnitCaster {
	uc.BeforeValue = g
	uc.AfterValue = g * 0.00220462262
	return uc
}

func (uc *UnitCaster) Float(precision ...int) float64 {
	p := uc.Precision
	if len(precision) != 0 {
		p = precision[0]
	}
	ratio := math.Pow(10, float64(p))
	return math.Round(uc.AfterValue*ratio) / ratio
}

func (uc *UnitCaster) Int() int {
	return int(uc.Float(0))
}
