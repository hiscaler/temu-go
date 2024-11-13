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
