package entity

import (
	"strings"
	"time"
)

// MallPermission 店铺权限
type MallPermission struct {
	MallId       int      `json:"mallId"`
	ExpiredTime  int64    `json:"expiredTime"`
	APIScopeList []string `json:"apiScopeList"`
}

// Accessible 是否有指定 API 权限
func (m MallPermission) Accessible(api string) bool {
	api = strings.TrimSpace(api)
	if api == "" {
		return false
	}

	for _, v := range m.APIScopeList {
		if strings.EqualFold(v, api) {
			return true
		}
	}
	return false
}

// Expired Token 是否过期
func (m MallPermission) Expired() bool {
	return m.ExpiredTime >= time.Now().UnixMilli()
}
