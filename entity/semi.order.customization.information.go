package entity

// SemiOrderCustomizationInformationPreview 半托订单定制信息预览
type SemiOrderCustomizationInformationPreview struct {
	CustomizedAreaId string `json:"customizedAreaId"` // 定制区域 ID. This field will only be returned when templateType=3, previewType=3 or 4
	ImageUrl         string `json:"imageUrl"`         // Image URL
	CustomizedText   string `json:"customizedText"`   // 定制文本
	// type of preview item, enum values:
	// - 1: overall preview image(If the product does not have a customized area configured, it represents the effect image uploaded by the merchant)
	// - 3: user uploaded image
	// - 4: customized text
	PreviewType int `json:"previewType"` // 预览类型
}

// SemiOrderCustomizationInformation 半托订单定制信息
type SemiOrderCustomizationInformation struct {
	// Customization template type when user created customized information, return null when there is no template for the product, enum values:
	// - 1: only image
	// - 2: only text
	// - 3: text and image
	TemplateType   int                                        `json:"templateType"`
	PreviewList    []SemiOrderCustomizationInformationPreview `json:"preview_list"`   // Graphic customization preview information, this field will only be returned when customizedType=2
	CustomizedData string                                     `json:"customizedData"` // Graphic customization content, in json format, this field will only be returned when customizedType=2
	OrderSn        string                                     `json:"orderSn"`        // OrderSn corresponding to customized information
	CustomizedText string                                     `json:"customizedText"` // Customization text, this field will only be returned when customizedType=1
	TemplateId     int                                        `json:"templateId"`     // Customization template ID when user created customized information, return null when there is no template for the product
	// Customized type, enum values:
	// - 1: pure text customization, no customized templates
	// - 2: customized graphics and text, with customized templates available
	CustomizedType int `json:"customizedType"`
}
