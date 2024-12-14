package normal

import (
	"github.com/hiscaler/temu-go/entity"
)

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
func (p *Pager) TidyPager(options ...int) *Pager {
	page, maxPageSize := 1, entity.MaxPageSize
	n := len(options)
	if n != 0 {
		page = options[0]
		if n >= 2 {
			maxPageSize = options[1]
		}
	}
	if p.Page <= 0 {
		p.Page = page
	}
	if p.PageSize <= 0 || p.PageSize > maxPageSize {
		p.PageSize = maxPageSize
	}
	return p
}
