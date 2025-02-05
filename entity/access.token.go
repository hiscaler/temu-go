package entity

type AccessToken struct {
	MallId       int64    `json:"mallId"`
	AccessToken  string   `json:"accessToken"`
	ExpiredTime  int64    `json:"expiredTime"`
	ApiScopeList []string `json:"apiScopeList"`
}
