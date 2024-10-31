package entity

import (
	"strings"
	"time"
)

// MallPermission 店铺权限
type MallPermission struct {
	MallID       int      `json:"mallId"`
	ExpiredTime  int64    `json:"expiredTime"`
	APIScopeList []string `json:"apiScopeList"`
}

// HasAPI 是否有指定 API 权限
func (m MallPermission) HasAPI(method string) bool {
	for _, v := range m.APIScopeList {
		if strings.EqualFold(v, method) {
			return true
		}
	}
	return false
}

// Expired Token 是否过期
func (m MallPermission) Expired() bool {
	return m.ExpiredTime >= time.Now().UnixMilli()
}
