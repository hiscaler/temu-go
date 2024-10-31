package entity

// GoodsSizeChartSetting 尺码分类对象
type GoodsSizeChartSetting struct {
	GroupChName    string   `json:"groupChName"`    // 尺码组中文名
	GroupEnName    string   `json:"groupEnName"`    // 尺码组英文名
	Code           int      `json:"code"`           // 尺码组 ID
	MappingContent struct{} `json:"mappingContent"` // 尺码参数，留意unnecessary=true，对应类目尺码组和尺码元素必填
}
