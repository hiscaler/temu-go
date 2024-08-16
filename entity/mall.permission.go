package entity

// MallPermission 店铺权限
type MallPermission struct {
	MallID       int      `json:"mallId"`
	ExpiredTime  int      `json:"expiredTime"`
	APIScopeList []string `json:"apiScopeList"`
}
