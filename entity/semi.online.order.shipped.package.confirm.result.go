package entity

type SemiOnlineOrderShippedPackageConfirmResult struct {
	PackageSn      string `json:"packageSn"`      // 包裹号
	WarningMessage string `json:"warningMessage"` // 警告消息
}
