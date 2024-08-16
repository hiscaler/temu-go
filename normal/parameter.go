package normal

type Parameter struct {
	AppKey      string `json:"app_key"`
	Timestamp   string `json:"timestamp"`
	Sign        string `json:"sign,omitempty"`
	DataType    string `json:"data_type"`
	AccessToken string `json:"access_token"`
	Version     string `json:"version,omitempty"`
}

type ParameterWithPager struct {
	Parameter
	Page     int `json:"pageNo"`
	PageSize int `json:"pageSize"`
}

func (pp *ParameterWithPager) TidyPager() *ParameterWithPager {
	if pp.Page <= 0 {
		pp.Page = 1
	}
	if pp.PageSize <= 0 {
		pp.PageSize = 10
	} else if pp.PageSize > 500 {
		pp.PageSize = 500
	}
	return pp
}
