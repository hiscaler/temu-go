package entity

// FreightTemplate 物流模板
type FreightTemplate struct {
	FreightTemplateId int64  `json:"freightTemplateId"` // 模板 id
	TemplateName      string `json:"templateName"`      // 模板名称
}

// SemiFreightTemplate 半托运费模板
type SemiFreightTemplate struct {
	FreightTemplateId string `json:"freightTemplateId"` // 模板 ID
	TemplateName      string `json:"templateName"`      // 模板名称
}
