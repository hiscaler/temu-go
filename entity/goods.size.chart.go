package entity

import "gopkg.in/guregu/null.v4"

// GoodsSizeChart 商品尺码
type GoodsSizeChart struct {
	BusinessId int       `json:"businessId"` // 模板ID
	Name       string    `json:"name"`       // 模板名称
	ClassId    int       `json:"classId"`    // 尺码分类ID
	Reusable   null.Bool `json:"reusable"`   // 是否可重复使用
	ContentVO  struct {
		Meta struct {
			Groups []struct {
				Id   int    `json:"id"`
				Name string `json:"name"`
			} `json:"groups"` // 尺码组元数据
			Elements []struct {
				Id   int    `json:"id"`
				Name string `json:"name"`
			} `json:"elements"` // 尺码参数组元数据
		} `json:"meta"` // 尺码组与尺码参数组元数据
		Records []map[int]string `json:"records"` // 元数据-值映射关系
	} `json:"contentVO"` // 内容
	UpdateAt int `json:"updateAt"` // 更新时间
}
