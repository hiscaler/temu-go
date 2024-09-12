package normal

type Parameter struct {
	AppKey      string `json:"app_key"`
	Timestamp   string `json:"timestamp"`
	Sign        string `json:"sign,omitempty"`
	DataType    string `json:"data_type"`
	AccessToken string `json:"access_token"`
	Version     string `json:"version,omitempty"`
}

type Pager struct {
	Page     int `json:"pageNo"`
	PageSize int `json:"pageSize"`
}

type ParameterWithPager struct {
	Parameter
	Pager
}

// TidyPager 设置翻页数据
func (pp *Pager) TidyPager(values ...int) *Pager {
	page := 1
	maxPageSize := 100
	n := len(values)
	if n != 0 {
		page = values[0]
		if n >= 2 {
			maxPageSize = values[1]
		}
	}
	if pp.Page <= 0 {
		pp.Page = page
	}
	if pp.PageSize <= 0 {
		pp.PageSize = 10
	} else if pp.PageSize > maxPageSize {
		pp.PageSize = maxPageSize
	}
	return pp
}
