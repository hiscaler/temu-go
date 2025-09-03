package entity

type SemiOrderLogisticsScanFormResult struct {
	ScanFormSn    string   `json:"scanFormSn"`    // 扫描单号
	PackageSnList []string `json:"packageSnList"` // 包裹单号列表
}
