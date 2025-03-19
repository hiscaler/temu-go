package entity

import "gopkg.in/guregu/null.v4"

// GoodsCategoryAttribute 分类属性
type GoodsCategoryAttribute struct {
	InputMaxSpecNum      int                              `json:"inputMaxSpecNum"`      // 模板允许的最大的自定义规格数量
	ChooseAllQualifySpec bool                             `json:"chooseAllQualifySpec"` // 限定规格是否必须全选
	SingleSpecValueNum   int                              `json:"singleSpecValueNum"`   // 单个自定义规格值上限
	Properties           []GoodsCategoryAttributeProperty `json:"properties"`           // 模板属性
}

// GoodsCategoryAttributeProperty 模板属性
type GoodsCategoryAttributeProperty struct {
	ParentSpecId                    null.Int    `json:"parentSpecId"`                    // 规格id
	NumberInputTitle                null.String `json:"numberInputTitle"`                // 数值录入Title
	TemplatePropertyValueParentList []any       `json:"templatePropertyValueParentList"` // 属性值关联关系
	MaxValue                        string      `json:"maxValue"`                        // 输入最大值：文本类型代表文本最长长度、 数值类型代表数字最大值、时间类型代表时间最大值
	Values                          []struct {
		SpecId        null.Int    `json:"specId"`
		Vid           int         `json:"vid"`
		Lang2Value    null.String `json:"lang2Value"`
		ParentVidList []int64     `json:"parentVidList"`
		ExtendInfo    null.String `json:"extendInfo"`
		Value         string      `json:"value"`
		Group         null.String `json:"group"`
	} `json:"values"` // 模板属性值
	ValueUnit      []string `json:"valueUnit"`      // 属性值单位
	ChooseMaxNum   int      `json:"chooseMaxNum"`   // 最大可勾选数目
	Pid            int64    `json:"pid"`            // 基础属性id
	TemplatePid    int64    `json:"templatePid"`    // 模板属性id
	Required       bool     `json:"required"`       // 属性是否必填
	InputMaxNum    int      `json:"inputMaxNum"`    // 最大可输入数目,为0时代表不可输入
	ValuePrecision int      `json:"valuePrecision"` // 小数点允许最大精度,为0时代表不允许输入小数
	// TEXT(0, "文本-传属性值id或者自定义的属性值"),
	// NUM(1, "数值-传属性值和单位"),
	// NUMBER_RANGE(2, "输入数值范围"),
	// NUMBER_PRODUCT_DOUBLE(3, "输入数值乘积-2维"),
	// NUMBER_PRODUCT_TRIPLE(4, "输入数值乘积-3维"),
	// SINGLE_YMD_DATE(5, "单项时间选择器-年月日"),
	// MULTIPLE_YMD_DATE(6, "双项时间选择器-年月日"),
	// SINGLE_YM_DATE(7, "单项时间选择器-年月"),
	// MULTIPLE_YM_DATE(8, "双项时间选择器-年月"),
	PropertyValueType int `json:"propertyValueType"` // 属性值类型
	// INPUT(0, "可输入"), CHOOSE(1, "可勾选"), INPUT_CHOOSE(3, "可输入又可勾选"), SINGLE_YMD_DATE(5, "单项时间选择器-年月日"), MULTIPLE_YMD_DATE(6, "双项时间选择器-年月日"), SINGLE_YM_DATE(7, "单项时间选择器-年月"), MULTIPLE_YM_DATE(8, "双项时间选择器-年月"), COLOR_SELECTOR(9, "调色盘"), SIZE_SELECTOR(10, "尺码选择器"), NUMBER_RANGE(11, "输入数值范围"), NUMBER_PRODUCT_DOUBLE(12, "输入数值乘积-2维"), NUMBER_PRODUCT_TRIPLE(13, "输入数值乘积-3维"), AUTO_COMPUTE(14, "自动计算框"), REGION_CHOOSE(15, "地区选择器"), PROPERTY_CHOOSE_AND_INPUT(16, "属性勾选和数值录入"),
	ControlType         int         `json:"controlType"`         // 控件类型
	ValueRule           int         `json:"valueRule"`           // 数值规则：SUM_OF_VALUES_IS_100(1, "数值之和等于100")
	PropertyChooseTitle null.String `json:"propertyChooseTitle"` // 属性勾选Title
	Name                string      `json:"name"`                // 属性名称
	Lang2Name           null.String `json:"lang2Name"`           // Lang2 Name
	IsSale              bool        `json:"isSale"`              // 是否销售属性(区分普通属性与规格属性)
	ShowType            int         `json:"showType"`            // B端展示规则
	ParentTemplatePid   int         `json:"parentTemplatePid"`   // 模板父属性ID
	MainSale            null.Int    `json:"mainSale"`            // 是否为主销售属性
	RefPid              int64       `json:"refPid"`              // 属性id
	// 如果子属性的控件类型是可勾选，则使用 templatePropertyValueParentList。如果子属性的控件类型是非可勾选，则用 showCondition
	ShowCondition []struct {
		ParentRefPid int64   `json:"parentRefPid"` // 父属性id
		ParentVids   []int64 `json:"parentVids"`   // 若属性按条件展示,则只有parent_vids中的值被选择时属性才可使用
	} `json:"showCondition"` // 属性展示条件，或者关系
}
