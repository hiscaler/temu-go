package entity

// Specification 规格
type Specification struct {
	SpecId         int    `json:"specId"`
	ParentSpecName string `json:"parentSpecName"`
	SpecName       string `json:"specName"`
	ParentSpecId   int    `json:"parentSpecId"`
}
