package temu

import (
	"context"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

type jitPresaleRuleService service

// Query jit预售规则查询接口（bg.virtualinventoryjit.rule.get）
// https://seller.kuajingmaihuo.com/sop/view/706628248275137588#9h0RVQ
// 全托管JIT开通：全托管的SKC开通JIT模式，需要签署对应协议之后才可添加虚拟库存
func (s jitPresaleRuleService) Query(ctx context.Context, productId, productSkcId int64) (rule entity.JitPresaleRule, err error) {
	var result = struct {
		normal.Response
		Result entity.JitPresaleRule `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(map[string]int64{"productId": productId, "productSkcId": productSkcId}).
		SetResult(&result).
		Post("bg.virtualinventoryjit.rule.get")
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	return result.Result, nil
}

// Sign jit预售规则签署接口（bg.virtualinventoryjit.rule.sign）
// https://seller.kuajingmaihuo.com/sop/view/706628248275137588#q8IeTi
// - 全托管JIT开通：全托管的SKC开通JIT模式，需要签署对应协议之后才可添加虚拟库存

type JitPresaleRuleSignRequest struct {
	ProductId      int64  `json:"productId"`      // 货品id，货品需要处于JIT开启状态，才能签署JIT协议
	AgtVersion     int    `json:"agtVersion"`     // JIT预售协议版本号
	ProductAgtType int    `json:"productAgtType"` // 货品协议类型，1-JIT模式快速售卖协议
	Url            string `json:"url"`            // JIT协议链接
}

func (m JitPresaleRuleSignRequest) validate() error {
	return nil
}

func (s jitPresaleRuleService) Sign(ctx context.Context, request JitPresaleRuleSignRequest) (ok bool, err error) {
	if err = request.validate(); err != nil {
		return
	}

	var result = struct {
		normal.Response
		Result any `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(request).
		SetResult(&result).
		Post("bg.virtualinventoryjit.rule.sign")
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	return true, nil
}
