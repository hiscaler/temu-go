package temu

import (
	"context"
	"crypto/md5"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	gci "github.com/echo-ok/goods-customization-information"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-resty/resty/v2"
	"github.com/hiscaler/gox/filex"
	"github.com/hiscaler/gox/randx"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/helpers"
	"github.com/hiscaler/temu-go/normal"
	"github.com/hiscaler/temu-go/validators/is"
	"gopkg.in/guregu/null.v4"
)

// 订单服务（半托管专属，必须在 US/EU 网关调用）
type semiOrderService service

type SemiOrderQueryParams struct {
	normal.ParameterWithPager
	PageNumber int `json:"pageNumber"` // 第几页
	// 父单状态，默认查全部枚举值如下
	// 0-全部
	// 1-”PENDING“，挂起中
	// 2-"UN_SHIPPING"，待发货
	// 3-"CANCELED",已取消
	// 4-"SHIPPED"，已发货
	// 5-“RECEIPTED”,已签收
	// 备注：
	// 本本订单还存在如下状态
	// 41-部分取消
	// 51-部分签收
	ParentOrderStatus         null.Int `json:"parentOrderStatus,omitempty"`         // 父单状态
	CreateBefore              string   `json:"createBefore,omitempty"`              // 父单创建时间结束查询时间，格式是秒时间戳。查询时间需要同时入参开始和结束时间才生效
	CreateAfter               string   `json:"createAfter,omitempty"`               // 父单创建时间开始查询时间，格式是秒时间戳
	ParentOrderSnList         []string `json:"parentOrderSnList,omitempty"`         // 父单号列表，单次请求最多 20 个
	ExpectShipLatestTimeStart string   `json:"expectShipLatestTimeStart,omitempty"` // 期望最晚发货时间开始查询时间，格式是秒时间戳
	ExpectShipLatestTimeEnd   string   `json:"expectShipLatestTimeEnd,omitempty"`   // 期望最晚发货时间结束查询时间，格式是秒时间戳。查询时间需要同时入参开始和结束时间才生效
	UpdateAtStart             string   `json:"updateAtStart,omitempty"`             // 订单更新时间开始查询时间，格式是秒时间戳
	UpdateAtEnd               string   `json:"updateAtEnd,omitempty"`               // 订单更新时间结束查询时间，格式是秒时间戳。查询时间需要同时入参开始和结束时间才生效
	RegionId                  int      `json:"regionId"`                            // 区域 ID，美国-211
	// 子订单履约类型，具体枚举值如下：
	// 1. 数组只传入 fulfillBySeller，只返回卖家履约子订单列表
	// 2. 数组只传入 fulfillByCooperativeWarehouse，只返回合作仓履约子订单列表
	// 3. 数组传入 fulfillBySeller 和 fulfillByCooperativeWarehouse，返回卖家履约子订单列表+合作仓履约子订单列表
	// 4. fulfillmentTypeList不传或者传了为空，默认返回卖家履约子订单列表
	FulfillmentTypeList []string    `json:"fulfillmentTypeList,omitempty"` // 子订单履约类型
	ParentOrderLabel    []string    `json:"parentOrderLabel,omitempty"`    // PO 单订单状态标签
	SortBy              null.String `json:"sortby,omitempty"`              // 排序依据，倒序输出。默认按照订单创建时间，对应枚举为：updateTime、createTime
}

func (m SemiOrderQueryParams) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.ParentOrderStatus, validation.When(m.ParentOrderStatus.Valid, validation.By(func(value interface{}) error {
			v, ok := value.(null.Int)
			if !ok {
				return errors.New("无效的父单状态")
			}

			return validation.Validate(int(v.Int64), validation.In(
				entity.SemiOrderStatusAll,
				entity.SemiOrderStatusPending,
				entity.SemiOrderStatusUnShipping,
				entity.SemiOrderStatusCanceled,
				entity.SemiOrderStatusShipped,
				entity.SemiOrderStatusReceipted,
				entity.SemiOrderStatusPartialCanceled,
				entity.SemiOrderStatusPartialReceipted,
			).Error("无效的父单状态"))
		}))),
		validation.Field(&m.CreateBefore,
			validation.When(m.CreateBefore != "" || m.CreateAfter != "", validation.By(is.TimeRange(m.CreateAfter, m.CreateBefore, time.DateTime))),
		),
		validation.Field(&m.ExpectShipLatestTimeStart,
			validation.When(m.ExpectShipLatestTimeStart != "" || m.ExpectShipLatestTimeEnd != "", validation.By(is.TimeRange(m.ExpectShipLatestTimeStart, m.ExpectShipLatestTimeEnd, time.DateTime))),
		),
		validation.Field(&m.UpdateAtStart,
			validation.When(m.UpdateAtStart != "" || m.UpdateAtEnd != "", validation.By(is.TimeRange(m.UpdateAtStart, m.UpdateAtEnd, time.DateTime))),
		),
		validation.Field(&m.RegionId, validation.By(is.RegionId(entity.RegionIds))),
		validation.Field(&m.FulfillmentTypeList, validation.Each(validation.By(func(value interface{}) error {
			v, ok := value.(string)
			if !ok {
				return errors.New("无效的子订单履约类型")
			}
			return validation.Validate(v, validation.In(
				entity.SemiOrderFulfillmentTypeBySeller,
				entity.SemiOrderFulfillmentTypeByCooperativeWarehouse,
			).Error("无效的子订单履约类型"))
		}))),
		validation.Field(&m.ParentOrderLabel, validation.Each(validation.By(func(value interface{}) error {
			v, ok := value.(string)
			if !ok {
				return errors.New("无效的父单状态标签")
			}
			return validation.Validate(v, validation.In(
				entity.SemiParentOrderLabelSoonToBeOverdue,
				entity.SemiParentOrderLabelPastDue,
				entity.SemiParentOrderLabelPendingBuyerCancellation,
				entity.SemiParentOrderLabelPendingBuyerAddressChange,
			).Error("无效的父单状态标签"))
		}))),
		validation.Field(&m.SortBy, validation.When(m.SortBy.Valid, validation.By(func(value interface{}) error {
			v, ok := value.(null.String)
			if !ok {
				return errors.New("无效的排序依据")
			}
			return validation.Validate(v.String, validation.In(
				entity.SemiOrderOrderByCreateTime,
				entity.SemiOrderOrderByUpdateTime,
			).Error("无效的排序依据"))
		}))),
	)
}

// Query 订单列表查询接口
// https://seller.kuajingmaihuo.com/sop/view/867739977041685428#r2WKrz
func (s semiOrderService) Query(ctx context.Context, params SemiOrderQueryParams) (items []entity.ParentOrder, total, totalPages int, isLastPage bool, err error) {
	params.PageNumber = params.TidyPager().Page
	params.OmitPage()
	if err = params.validate(); err != nil {
		err = invalidInput(err)
		return
	}

	if params.CreateBefore != "" && params.CreateAfter != "" {
		if start, end, e := helpers.StrTime2Unix(params.CreateAfter, params.CreateBefore); e == nil {
			params.CreateBefore = end
			params.CreateAfter = start
		}
	}

	if params.ExpectShipLatestTimeStart != "" && params.ExpectShipLatestTimeEnd != "" {
		if start, end, e := helpers.StrTime2Unix(params.ExpectShipLatestTimeStart, params.ExpectShipLatestTimeEnd); e == nil {
			params.ExpectShipLatestTimeStart = start
			params.ExpectShipLatestTimeEnd = end
		}
	}

	if params.UpdateAtStart != "" && params.UpdateAtEnd != "" {
		if start, end, e := helpers.StrTime2Unix(params.UpdateAtStart, params.UpdateAtEnd); e == nil {
			params.UpdateAtStart = start
			params.UpdateAtEnd = end
		}
	}

	var result = struct {
		normal.Response
		Result struct {
			TotalItemNum int                  `json:"totalItemNum"`
			PageItems    []entity.ParentOrder `json:"pageItems"`
		} `json:"result"`
	}{}

	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(params).
		SetResult(&result).
		Post("bg.order.list.v2.get")
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	items = result.Result.PageItems
	total, totalPages, isLastPage = parseResponseTotal(params.Page, params.PageSize, result.Result.TotalItemNum)
	return
}

// ShippingInformation 订单收货地址查询接口（bg.order.shippinginfo.v2.get）
// https://seller.kuajingmaihuo.com/sop/view/867739977041685428#AVEKr6
func (s semiOrderService) ShippingInformation(ctx context.Context, parentOrderNumber string) (entity.SemiOrderShippingInformation, error) {
	var result = struct {
		normal.Response
		Result entity.SemiOrderShippingInformation `json:"result"`
	}{}

	resp, err := s.httpClient.
		R().
		SetContext(ctx).
		SetBody(map[string]string{"parentOrderSn": parentOrderNumber}).
		SetResult(&result).
		Post("bg.order.shippinginfo.v2.get")
	if err = recheckError(resp, result.Response, err); err != nil {
		return entity.SemiOrderShippingInformation{}, err
	}

	return result.Result, nil
}

// CustomizationInformation 半托订单定制信息查询
// https://partner.temu.com/documentation?menu_code=fb16b05f7a904765aac4af3a24b87d4a&sub_menu_code=e8f86a2f5241441e9b095bf309d04dce
// 注意：orderNumbers 为子单号
func (s semiOrderService) CustomizationInformation(ctx context.Context, orderNumbers ...string) ([]entity.SemiOrderCustomizationInformation, error) {
	if len(orderNumbers) == 0 {
		return nil, errors.New("待查询半托订单定制信息子单号列表不能为空")
	}

	var result = struct {
		normal.Response
		Result []entity.SemiOrderCustomizationInformation `json:"result"`
	}{}
	resp, err := s.httpClient.
		R().
		SetContext(ctx).
		SetBody(map[string][]string{"orderSnList": orderNumbers}).
		SetResult(&result).
		Post("bg.order.customization.get")
	if err = recheckError(resp, result.Response, err); err != nil {
		return nil, err
	}

	results := result.Result
	if len(results) == 0 {
		return results, nil
	}

	imagePreviews := make(map[int]map[int]entity.SemiOrderCustomizationInformationPreview)
	for i, res := range results {
		for j, v := range res.PreviewList {
			if v.PreviewType == 1 || v.PreviewType == 3 {
				if _, ok := imagePreviews[i]; !ok {
					imagePreviews[i] = make(map[int]entity.SemiOrderCustomizationInformationPreview)
				}
				imagePreviews[i][j] = v
			}
		}
	}
	if len(imagePreviews) == 0 {
		return results, nil
	}

	keys := []string{
		"toa-access-token",
		"toa-app-key",
		"toa-random",
		"toa-timestamp",
	}
	dir := "./static_files/temu/semi/images"
	sb := strings.Builder{}
	headers := map[string]string{
		"toa-app-key":      s.config.AppKey,
		"toa-access-token": s.config.AccessToken,
	}
	httpClient := resty.New().
		SetDebug(s.debug).
		SetHeaders(map[string]string{
			"Content-Type": "application/json",
			"User-Agent":   UserAgent,
		}).
		SetAllowGetMethodPayload(true).
		SetTimeout(s.config.Timeout * time.Second).
		SetTransport(&http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: !s.config.VerifySSL},
			DialContext: (&net.Dialer{
				Timeout: s.config.Timeout * time.Second,
			}).DialContext,
		})
	if s.debug {
		httpClient.EnableTrace()
	}
	for i, res := range imagePreviews {
		for j, pl := range res {
			if !pl.ImageUrl.Valid {
				continue
			}

			var parsedURL *url.URL
			parsedURL, err = url.Parse(pl.ImageUrl.String)
			if err != nil {
				results[i].PreviewList[j].ImageDownloadError = null.StringFrom(err.Error())
				continue
			}
			filename := strings.ToLower(filepath.Base(parsedURL.Path))
			if filename == "" {
				results[i].PreviewList[j].ImageDownloadError = null.StringFrom("无法获取文件名")
			}
			savePath := filepath.Join(dir, filename)
			if filex.Exists(savePath) {
				results[i].PreviewList[j].ImageDownloadUrl = null.StringFrom(urlJoin(s.config.StaticFileServer, savePath))
				continue
			}

			headers["toa-random"] = randx.Letter(32, true)
			headers["toa-timestamp"] = strconv.FormatInt(time.Now().Unix(), 10)
			sb.Reset()
			sb.WriteString(s.config.AppSecret)
			for _, key := range keys {
				sb.WriteString(key)
				sb.WriteString(headers[key])
			}
			sb.WriteString(s.config.AppSecret)
			headers["toa-sign"] = strings.ToUpper(fmt.Sprintf("%x", md5.Sum([]byte(sb.String()))))
			resp, err = httpClient.
				SetOutputDirectory(dir).
				R().
				SetHeaders(headers).
				SetOutput(filename).
				Get(pl.ImageUrl.String)
			if err != nil {
				results[i].PreviewList[j].ImageDownloadError = null.StringFrom(err.Error())
			} else {
				if resp.IsError() {
					results[i].PreviewList[j].ImageDownloadError = null.StringFrom(resp.String())
				} else if resp.IsSuccess() {
					results[i].PreviewList[j].ImageDownloadUrl = null.StringFrom(urlJoin(s.config.StaticFileServer, path.Join(dir, filename)))
				} else {
					results[i].PreviewList[j].ImageDownloadError = null.StringFrom(resp.String())
				}
			}
		}
	}

	return results, nil
}

// CustomizationInformation2 半托订单定制信息查询
// https://partner.temu.com/documentation?menu_code=fb16b05f7a904765aac4af3a24b87d4a&sub_menu_code=e8f86a2f5241441e9b095bf309d04dce
// 注意：orderNumbers 为子单号
func (s semiOrderService) CustomizationInformation2(ctx context.Context, orderNumbers ...string) ([]gci.GoodsCustomizedInformation, error) {
	if len(orderNumbers) == 0 {
		return nil, errors.New("待查询半托订单定制信息子单号列表不能为空")
	}

	var result = struct {
		normal.Response
		Result []entity.SemiOrderCustomizationInformation `json:"result"`
	}{}
	resp, err := s.httpClient.
		R().
		SetContext(ctx).
		SetBody(map[string][]string{"orderSnList": orderNumbers}).
		SetResult(&result).
		Post("bg.order.customization.get")
	if err = recheckError(resp, result.Response, err); err != nil {
		return nil, err
	}

	results := result.Result
	if len(results) == 0 {
		return nil, nil
	}

	ciList := make([]gci.GoodsCustomizedInformation, 0)
	for _, res := range results {
		ci := gci.NewGoodsCustomizedInformation()
		ci.SetRawData(res)
		surface := gci.NewSurface()
		for _, v := range res.PreviewList {
			region := gci.NewRegion()
			var img gci.Image
			if v.PreviewType == 1 {
				if img, err = gci.NewImage(v.ImageUrl.ValueOrZero(), false); err == nil {
					surface.PreviewImage = &img
				}
			} else if v.PreviewType == 3 {
				if img, err = gci.NewImage(v.ImageUrl.ValueOrZero(), true); err == nil {
					region.AddImage(img)
				}
			}
			if v.CustomizedText.Valid {
				var text gci.Text
				if text, err = gci.NewText("", v.CustomizedText.String); err == nil {
					region.AddText(text)
				}
			}
			surface.AddRegion(region)
		}
		ci.AddSurface(surface)
		ciList = append(ciList, ci)
	}
	if len(ciList) == 0 {
		return ciList, nil
	}

	keys := []string{
		"toa-access-token",
		"toa-app-key",
		"toa-random",
		"toa-timestamp",
	}
	dir := "./static_files/temu/semi/images"
	sb := strings.Builder{}
	headers := map[string]string{
		"toa-app-key":      s.config.AppKey,
		"toa-access-token": s.config.AccessToken,
	}
	httpClient := resty.New().
		SetDebug(s.debug).
		SetHeaders(map[string]string{
			"Content-Type": "application/json",
			"User-Agent":   UserAgent,
		}).
		SetAllowGetMethodPayload(true).
		SetTimeout(s.config.Timeout * time.Second).
		SetTransport(&http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: !s.config.VerifySSL},
			DialContext: (&net.Dialer{
				Timeout: s.config.Timeout * time.Second,
			}).DialContext,
		})
	if s.debug {
		httpClient.EnableTrace()
	}
	for i, ci := range ciList {
		for j, surface := range ci.Surfaces {
			if surface.PreviewImage == nil || !surface.PreviewImage.Redownload() {
				continue
			}

			var parsedURL *url.URL
			parsedURL, err = url.Parse(surface.PreviewImage.RawUrl)
			if err != nil {
				ciList[i].Surfaces[j].PreviewImage.SetError(err)
				continue
			}
			filename := strings.ToLower(filepath.Base(parsedURL.Path))
			if filename == "" {
				ciList[i].Surfaces[j].PreviewImage.SetError("无法获取文件名")
			}
			savePath := filepath.Join(dir, filename)
			if filex.Exists(savePath) {
				ciList[i].Surfaces[j].PreviewImage.SetUrl(urlJoin(s.config.StaticFileServer, savePath))
				continue
			}

			headers["toa-random"] = randx.Letter(32, true)
			headers["toa-timestamp"] = strconv.FormatInt(time.Now().Unix(), 10)
			sb.Reset()
			sb.WriteString(s.config.AppSecret)
			for _, key := range keys {
				sb.WriteString(key)
				sb.WriteString(headers[key])
			}
			sb.WriteString(s.config.AppSecret)
			headers["toa-sign"] = strings.ToUpper(fmt.Sprintf("%x", md5.Sum([]byte(sb.String()))))
			resp, err = httpClient.
				SetOutputDirectory(dir).
				R().
				SetHeaders(headers).
				SetOutput(filename).
				Get(surface.PreviewImage.RawUrl)
			if err != nil {
				ciList[i].Surfaces[j].PreviewImage.SetError(err)
			} else {
				if resp.IsError() {
					ciList[i].Surfaces[j].PreviewImage.SetError(resp.String())
				} else if resp.IsSuccess() {
					ciList[i].Surfaces[j].PreviewImage.SetUrl(urlJoin(s.config.StaticFileServer, path.Join(dir, filename)))
				} else {
					ciList[i].Surfaces[j].PreviewImage.SetError(resp.String())
				}
			}
		}
	}

	return ciList, nil
}
