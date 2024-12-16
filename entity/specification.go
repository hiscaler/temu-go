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
