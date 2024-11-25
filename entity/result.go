package entity

import "gopkg.in/guregu/null.v4"

// Result 处理结果
type Result struct {
	Key     string      `json:"key"`
	Success bool        `json:"success"`
	Code    null.Int    `json:"code"`
	Error   null.String `json:"error"`
}
