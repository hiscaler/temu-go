package temu

import (
	"crypto/md5"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/hiscaler/gox/stringx"
	"github.com/hiscaler/temu-go/config"
	"github.com/hiscaler/temu-go/normal"
	"log"
	"net"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

const (
	Version   = "0.0.1"
	UserAgent = "Temu API Client-Golang/" + Version
)

const (
	BadRequestError           = 400    // 错误的请求
	UnauthorizedError         = 401    // 身份验证或权限错误
	NotFoundError             = 404    // 访问资源不存在
	InternalServerError       = 500    // 服务器内部错误
	MethodNotImplementedError = 501    // 方法未实现
	SystemExceptionError      = 200000 // 系统异常
)

var ErrNotFound = errors.New("not found")
var ErrInvalidParameters = errors.New("invalid parameters")

func IntPtr(value int) *int {
	return &value
}

type service struct {
	debug      bool          // Is debug mode
	logger     *log.Logger   // Log
	httpClient *resty.Client // HTTP client
}

type services struct {
	ShipOrder               shipOrderService
	ShipOrderStaging        shipOrderStagingService
	ShipOrderPacking        shipOrderPackingService
	ShipOrderPackage        shipOrderPackageService
	Barcode                 barcodeService
	PurchaseOrder           purchaseOrderService
	GoodsSales              goodsSalesService
	Logistics               logisticsService
	VirtualInventoryJit     virtualInventoryJitService
	GoodsSizeChart          goodsSizeChartService
	GoodsSizeChartClass     goodsSizeChartClassService
	GoodsSizeChartSetting   goodsSizeChartSettingService
	MallAddress             mallAddressService
	ShipOrderReceiveAddress shipOrderReceiveAddressService
	Goods                   goodsService
	Mall                    mallService
}

type Client struct {
	Debug        bool           // Is debug mode
	Logger       *log.Logger    // Log
	Services     services       // API services
	TimeLocation *time.Location // 时区
}

func NewClient(config config.Config) *Client {
	logger := log.New(os.Stdout, "[ Temu ] ", log.LstdFlags|log.Llongfile)
	client := &Client{
		Debug:  config.Debug,
		Logger: logger,
	}
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		loc = time.FixedZone("CST", 8*3600)
	}
	client.TimeLocation = loc

	httpClient := resty.New().
		SetDebug(config.Debug).
		EnableTrace().
		SetBaseURL("https://openapi.kuajingmaihuo.com/openapi/router").
		SetHeaders(map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
			"User-Agent":   UserAgent,
		}).
		SetAllowGetMethodPayload(true).
		SetTimeout(config.Timeout * time.Second).
		SetTransport(&http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: !config.VerifySSL},
			DialContext: (&net.Dialer{
				Timeout: config.Timeout * time.Second,
			}).DialContext,
		}).
		OnBeforeRequest(func(client *resty.Client, request *resty.Request) error {
			values := make(map[string]any)
			if request.Body != nil {
				b, err := json.Marshal(request.Body)
				if err != nil {
					return err
				}

				err = json.Unmarshal(b, &values)
				if err != nil {
					return err
				}
			}
			values["app_key"] = config.AppKey
			values["app_secret"] = config.AppSecret
			values["access_token"] = config.AccessToken
			values["data_type"] = "JSON"
			values["version"] = "V1"
			values["timestamp"] = time.Now().Unix()
			values["type"] = request.URL
			keys := make([]string, len(values))
			i := 0
			for k := range values {
				keys[i] = k
				i++
			}
			sort.Slice(keys, func(i, j int) bool {
				return keys[i] < keys[j]
			})
			stringBuilder := strings.Builder{}
			stringBuilder.WriteString(config.AppSecret)
			for _, key := range keys {
				stringBuilder.WriteString(key)
				stringBuilder.WriteString(stringx.String(values[key]))
			}
			stringBuilder.WriteString(config.AppSecret)
			values["sign"] = strings.ToUpper(fmt.Sprintf("%x", md5.Sum([]byte(stringBuilder.String()))))
			request.URL = ""
			request.SetBody(values)
			return nil
		}).
		SetRetryCount(3).
		SetRetryWaitTime(time.Duration(2) * time.Second).
		SetRetryMaxWaitTime(time.Duration(6) * time.Second).
		AddRetryCondition(func(response *resty.Response, err error) bool {
			if response == nil {
				return false
			}

			retry := response.StatusCode() == http.StatusTooManyRequests
			if !retry {
				r := struct {
					Success   bool   `json:"success"`
					ErrorCode int    `json:"errorCode"`
					ErrorMsg  string `json:"errorMsg"`
				}{}
				retry = json.Unmarshal(response.Body(), &r) == nil &&
					r.Success == false &&
					r.ErrorCode == 4000000 &&
					strings.EqualFold(r.ErrorMsg, "SYSTEM_EXCEPTION")
			}
			if retry {
				text := response.Request.URL
				if err != nil {
					text += fmt.Sprintf(", error: %s", err.Error())
				}
				logger.Printf("Retry request: %s", text)
			}
			return retry
		})
	if config.Debug {
		httpClient.SetBaseURL("https://openapi.kuajingmaihuo.com/openapi/router")
	}
	xService := service{
		debug:      config.Debug,
		logger:     logger,
		httpClient: httpClient,
	}
	client.Services = services{
		ShipOrder:               (shipOrderService)(xService),
		ShipOrderStaging:        (shipOrderStagingService)(xService),
		ShipOrderPacking:        (shipOrderPackingService)(xService),
		ShipOrderPackage:        (shipOrderPackageService)(xService),
		Barcode:                 (barcodeService)(xService),
		PurchaseOrder:           (purchaseOrderService)(xService),
		GoodsSales:              (goodsSalesService)(xService),
		Logistics:               (logisticsService)(xService),
		VirtualInventoryJit:     (virtualInventoryJitService)(xService),
		GoodsSizeChart:          (goodsSizeChartService)(xService),
		GoodsSizeChartClass:     (goodsSizeChartClassService)(xService),
		GoodsSizeChartSetting:   (goodsSizeChartSettingService)(xService),
		MallAddress:             (mallAddressService)(xService),
		ShipOrderReceiveAddress: (shipOrderReceiveAddressService)(xService),
		Goods:                   (goodsService)(xService),
		Mall:                    (mallService)(xService),
	}

	return client
}

func parseResponseTotal(currentPage, pageSize, total int) (n, totalPages int, isLastPage bool) {
	if currentPage == 0 {
		currentPage = 1
	}

	return total, (total / pageSize) + 1, currentPage >= totalPages
}

func parseResponse(resp *resty.Response, result normal.Response) (err error) {
	resp.Result()
	if resp.IsError() {
		errorMessage := result.ErrorMessage
		if errorMessage == "" {
			return errorWrap(resp.StatusCode(), resp.Error().(string))
		}
		return errors.New(errorMessage)
	}

	if !result.Success {
		return errors.New(result.ErrorMessage)
	}
	return nil
}

// errorWrap wrap an error, if status code is 200, return nil, otherwise return an error
func errorWrap(code int, message string) error {
	if code == http.StatusOK {
		return nil
	}

	if code == NotFoundError {
		return ErrNotFound
	}

	message = strings.TrimSpace(message)
	if message == "" {
		switch code {
		case BadRequestError:
			message = "请求错误"
		case UnauthorizedError:
			message = "认证失败，请确认您是否有相应的权限"
		case InternalServerError:
			message = "服务器内容错误"
		case MethodNotImplementedError:
			message = "方法未实现"
		case SystemExceptionError:
			message = "系统异常"
		default:
			message = "未知错误"
		}
	}
	return fmt.Errorf("%d: %s", code, message)
}
