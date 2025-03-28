package entity

import (
	"slices"
	"strings"
	"time"
)

// MallPermission 店铺权限
type MallPermission struct {
	MallId       int      `json:"mallId"`       // TEMU 店铺 ID
	ExpiredTime  int64    `json:"expiredTime"`  // 过期时间（时间戳秒级）
	APIScopeList []string `json:"apiScopeList"` // 有权限的 API 列表
}

// Accessible 是否有指定 API 权限
func (m MallPermission) Accessible(api string) bool {
	api = strings.TrimSpace(api)
	if api == "" {
		return false
	}

	return slices.ContainsFunc(m.APIScopeList, func(v string) bool {
		return strings.EqualFold(api, v)
	})
}

// Valid Token 是否有效
// days 如果大于零，则表示几天后过期，为零或者小于零则表示当前是否过期
func (m MallPermission) Valid(days ...int) bool {
	now := time.Now()
	if len(days) != 0 && days[0] > 0 {
		now = now.AddDate(0, 0, -days[0])
	}
	return m.ExpiredTime > now.Unix()
}
