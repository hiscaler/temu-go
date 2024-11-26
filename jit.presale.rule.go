package temu

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

type jitPresaleRuleService service

// Query jit预售规则查询接口（bg.virtualinventoryjit.rule.get）
// https://seller.kuajingmaihuo.com/sop/view/706628248275137588#9h0RVQ
// 全托管JIT开通：全托管的SKC开通JIT模式，需要签署对应协议之后才可添加虚拟库存
func (s jitPresaleRuleService) Query(ctx context.Context) (rule entity.JitPresaleRule, err error) {
	var result = struct {
		normal.Response
		Result entity.JitPresaleRule `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
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
	ProductId      int64  `json:"productId"`      // 货品 id，货品需要处于 JI T开启状态，才能签署 JIT 协议
	AgtVersion     int    `json:"agtVersion"`     // JIT 预售协议版本号
	ProductAgtType int    `json:"productAgtType"` // 货品协议类型（1: JIT模式快速售卖协议）
	Url            string `json:"url"`            // JIT 协议链接
}

func (m JitPresaleRuleSignRequest) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.ProductId,
			validation.Required.Error("无效的货品"),
		),
		validation.Field(&m.AgtVersion,
			validation.Required.Error("无效的 JIT 预售协议版本号"),
		),
		validation.Field(&m.ProductAgtType,
			validation.In(1).Error("无效的货品协议类型"),
		),
		validation.Field(&m.Url,
			validation.Required.Error("JIT 协议链接不能为空"),
			is.URL.Error("无效的 JIT 协议链接"),
		),
	)
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
