package entity

// FreightTemplate 物流模版
type FreightTemplate struct {
	FreightTemplateId int64  `json:"freightTemplateId"` // 模板 id
	TemplateName      string `json:"templateName"`      // 模板名称
}
