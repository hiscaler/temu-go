package entity

// GoodsCertification 商品资质
type GoodsCertification struct {
	CertType     int    `json:"certType"`     // 资质类型
	RejectReason string `json:"rejectReason"` // 驳回原因
	CertName     string `json:"certName"`     // 资质名称
	UpdateStatus int    `json:"updateStatus"` // 更新状态
	UpdateReason string `json:"updateReason"` // 更新原因
	AuditStatus  int    `json:"auditStatus"`  // 审核状态
}

// GoodsCertificationNeedUploadItem 资质要上传的内容
type GoodsCertificationNeedUploadItem struct {
	AliasName string `json:"aliasName"` // 文件名称
	// 当前资质类型需上传文件类型，枚举1-4，明细如下。
	// 返回多个时，”提交资质内容“时相关字段均需要上传。
	// - contentType：1  含义：资质证书
	//  - 必传字段：productCertFiles-货品资质文件列表
	// - contentType：2  含义：检测报告
	//  - 必传字段：inspectReportFiles-检测报告文件列表
	// - contentType：3  含义：资质编号
	//  - 必传字段：authCode-资质编号
	// - contentType：4  含义：产品实拍图
	//  - 使用“bg.flash.open.upload.real.image”上传
	//  - 必传字段：realPictures-实物图列表
	ContentType int `json:"contentType"` // 文件类型
	// hasExpireTime = true时，提交接口中的expireTime需不为null
	HasExpireTime bool `json:"hasExpireTime"` // 资质是否存在有效期
	ExpireDays    int  `json:"expireDays"`    // 固定失效天数
	// -needShowCustomer返回“true”时
	// bg.arbok.open.cert.uploadProductCert，需要入参showCustomer
	NeedShowCustomer bool `json:"needShowCustomer"` // 是否需要c端展示
	// 文件失效形式
	// 对不同类型资质的不同类型文件采取不同的管控方式，有固定日期，固定时长，不固定三种
	// - 1：不固定方式：失效日期不为null
	// - 2：固定日期方式：失效日期不为null
	// - 3：固定时长方式：失效日期不为null且生效日期不为null
	ExpireTime       int `json:"expireTime"`       // 固定失效日期
	ExpireNoticeDays int `json:"expireNoticeDays"` // 提前预警天数
}
