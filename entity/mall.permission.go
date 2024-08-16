package entity

// MallPermission 店铺权限
type MallPermission struct {
	MallID       int64    `json:"mallId"`
	ExpiredTime  int      `json:"expiredTime"`
	APIScopeList []string `json:"apiScopeList"`
}
