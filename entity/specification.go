package entity

import "strings"

// Specification 规格
type Specification struct {
	ParentSpecId   int    `json:"parentSpecId"`   // 父规格 ID
	ParentSpecName string `json:"parentSpecName"` // 父规格名称
	SpecId         int    `json:"specId"`         // 子规格 id
	SpecName       string `json:"specName"`       // 子规格名称
}

func (spec Specification) Ids() []int {
	rawIds := []int{spec.ParentSpecId, spec.SpecId}
	ids := make([]int, 0, len(rawIds))
	for _, id := range rawIds {
		if id == 0 {
			continue
		}
		ids = append(ids, id)
	}
	return ids
}

func (spec Specification) Names() []string {
	rawNames := []string{spec.ParentSpecName, spec.SpecName}
	names := make([]string, 0, len(rawNames))
	for _, name := range rawNames {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		names = append(names, name)
	}
	return names
}
