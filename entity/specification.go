package entity

import "strings"

// Specification 规格
type Specification struct {
	SpecId         int    `json:"specId"`
	SpecName       string `json:"specName"`
	ParentSpecId   int    `json:"parentSpecId"`
	ParentSpecName string `json:"parentSpecName"`
}

func (spec Specification) Ids() []int {
	ids := make([]int, 0)
	if spec.ParentSpecId != 0 {
		ids = append(ids, spec.ParentSpecId)
	}
	if spec.SpecId != 0 {
		ids = append(ids, spec.SpecId)
	}
	return ids
}

func (spec Specification) Names() []string {
	names := make([]string, 0)
	name := strings.TrimSpace(spec.ParentSpecName)
	if name != "" {
		names = append(names, name)
	}
	name = strings.TrimSpace(spec.SpecName)
	if name != "" {
		names = append(names, name)
	}
	return names
}
