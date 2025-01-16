package entity

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

// SemiOnlineOrderLogisticsChannel 物流供应商渠道
type SemiOnlineOrderLogisticsChannel struct {
	ChannelId             int64  `json:"channelId"`             // 渠道 id
	ShipLogisticsType     string `json:"shipLogisticsType"`     // 物流产品类型
	ShipCompanyId         int64  `json:"shipCompanyId"`         // 物流公司 id
	ShippingCompanyName   string `json:"shippingCompanyName"`   // 物流公司名称
	EstimatedText         string `json:"estimatedText"`         // 预估参数文案，包含预估的面单价格，币种，时效等信息 如：$39.46;USD ;1-2工作日送达
	EstimatedAmount       string `json:"estimatedAmount"`       // 预估金额，如：12.12
	EstimatedCurrencyCode string `json:"estimatedCurrencyCode"` // 预估货币币种
	InfoNeeded            []any  `json:"infoNeeded"`            // 渠道提示信息，提示这个渠道下 call 时候的依赖项目
	SignServiceName       string `json:"signServiceName"`       // 签收服务类型
	SignServiceId         string `json:"signServiceId"`         // 签收服务 ID
}

// Amount 获取预估金额
func (c SemiOnlineOrderLogisticsChannel) Amount() (float64, error) {
	s := strings.TrimSpace(c.EstimatedAmount)
	if s == "" {
		return 0, errors.New("预估金额待解析文本不能为空")
	}

	re, err := regexp.Compile(`\.?\d+\.?\d+`)
	if err != nil {
		return 0, err
	}

	v, err := strconv.ParseFloat(re.FindString(c.EstimatedAmount), 64)
	if err != nil {
		return 0, err
	}
	if v <= 0 {
		return 0, errors.New("无效的预估金额")
	}
	return v, nil
}

// DeliveryDays 获取交货天数
func (c SemiOnlineOrderLogisticsChannel) DeliveryDays() (int, int, error) {
	s := strings.TrimSpace(c.EstimatedText)
	if s == "" {
		return 0, 0, errors.New("交货天数待解析文本不能为空")
	}

	re, err := regexp.Compile(`([0-9]+)\s*-\s*([0-9]+)`)
	if err != nil {
		return 0, 0, err
	}

	values := re.FindStringSubmatch(s)
	if len(values) != 3 {
		return 0, 0, errors.New("交货天数文本解析失败")
	}

	var v, minDays, maxDays int
	if v, err = strconv.Atoi(values[1]); err != nil {
		return 0, 0, err
	} else {
		minDays = v
	}
	if v, err = strconv.Atoi(values[2]); err != nil {
		return 0, 0, err
	} else {
		maxDays = v
	}
	if minDays <= 0 || maxDays <= 0 {
		return 0, 0, errors.New("无效的交货天数")
	}

	return minDays, maxDays, nil
}
