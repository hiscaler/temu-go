package entity

import (
	"crypto/md5"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/hiscaler/gox/randx"
	"github.com/hiscaler/temu-go/config"
)

// SignatureUrl
// 加签文件处理逻辑
// 使用如下方式访问
// 访问地址：这里返回的 url 10 分钟后失效，超时请重新调用相应接口获得新的 url
// 请求方式：GET
// 请求 header 包含以下 5 个参数：
//
//	{
//		"toa-app-key": ${app_key},
//		"toa-access-token": ${access_token},
//		"toa-random": "${32 random numbers}",
//		"toa-timestamp":"${timestamp}",
//		"toa-sign":"${sign}"
//	}
//
// 其中 toa-random 自定义 32 个随机数
// 其中 toa-timestamp 长度要求 10 位，精确到秒级，current_time-300<=toa-timestamp<=current_time+300
// 其中 toa-sign 的计算逻辑与正常请求里公共参数 sign 的计算逻辑一致：
// - 将 toa-access-token 和 toa-app-key 和 toa-random 和 toa-timestamp 4 个参数进行首字母以 ASCII 方式升序排列 ascii asc，对于相同字母则使用下个字母做二次排序，字母序为从左到右，以此类推
// - 排序后的结果按照参数名 $key 参数值 $value 的次序进行字符串拼接，拼接处不包含任何字符
// - 拼接完成的字符串做进一步拼接成 1 个字符串（包含所有kv字符串的长串），并在该长串的头部及尾部分别拼接 app_secret，完成签名字符串的组装
// - 最后对签名字符串，使用 MD5 算法加密后，得到的 MD5 加密密文后转为大写，即为 toa-sign 值
// 接口返回文件流。
type SignatureUrl string

type File struct {
	Filename string `json:"filename"` // 文件名
	Content  []byte `json:"content"`  // 文件内容
}

// Decode 解码加签 URL，返回文件名，内容以及解码过程中出现的错误
func (su SignatureUrl) Decode(cfg config.Config) (f File, err error) {
	rawUrl := string(su)
	if rawUrl == "" {
		return f, errors.New("url is empty")
	}
	keys := []string{
		"toa-access-token",
		"toa-app-key",
		"toa-random",
		"toa-timestamp",
	}
	sb := strings.Builder{}
	headers := map[string]string{
		"toa-app-key":      cfg.AppKey,
		"toa-access-token": cfg.AccessToken,
	}
	httpClient := resty.New().
		SetDebug(cfg.Debug).
		SetHeaders(map[string]string{
			"Content-Type": "application/json",
		}).
		SetAllowGetMethodPayload(true).
		SetTimeout(cfg.Timeout * time.Second).
		SetTransport(&http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: !cfg.VerifySSL},
			DialContext: (&net.Dialer{
				Timeout: cfg.Timeout * time.Second,
			}).DialContext,
		})
	if cfg.Debug {
		httpClient.EnableTrace()
	}

	var u *url.URL
	u, err = url.Parse(rawUrl)
	if err != nil {
		return f, err
	}
	filename := strings.ToLower(filepath.Base(u.Path))
	if filename == "" {
		return f, errors.New("无法获取文件名")
	}
	f.Filename = filename

	headers["toa-random"] = randx.Letter(32, true)
	headers["toa-timestamp"] = strconv.FormatInt(time.Now().Unix(), 10)
	sb.Reset()
	sb.WriteString(cfg.AppSecret)
	for _, key := range keys {
		sb.WriteString(key)
		sb.WriteString(headers[key])
	}
	sb.WriteString(cfg.AppSecret)
	headers["toa-sign"] = strings.ToUpper(fmt.Sprintf("%x", md5.Sum([]byte(sb.String()))))
	resp, err := httpClient.R().SetHeaders(headers).SetOutput(filename).Get(rawUrl)
	if err != nil {
		return f, err
	}
	if !resp.IsSuccess() {
		return f, errors.New(resp.String())
	}
	f.Content = resp.Body()
	return f, nil
}
