package config

import (
	"log/slog"

	restyslog "github.com/hiscaler/resty-slog"
)

type EnvUrl struct {
	Prod string
	Test string
}

type RegionEnvUrls = map[string]EnvUrl

// LogMask describes the calling application's HTTP log masking policy.
type LogMask = restyslog.LogMask

type Config struct {
	Env              string             `json:"env"`                // 环境（dev：开发环境、test：测试环境、prod：生产环境）
	Debug            bool               `json:"debug"`              // 是否为调试模式
	LogLevel         string             `json:"log_level"`          // 调用日志级别，空值兼容 Debug
	LogMask          *restyslog.LogMask `json:"log_mask,omitempty"` // 调用方指定的请求及响应日志脱敏规则
	Region           string             `json:"region"`             // 区域
	AppKey           string             `json:"app_key"`            // App Key
	AppSecret        string             `json:"app_secret"`         // App 秘钥
	AccessToken      string             `json:"access_token"`       // Access Token
	Timeout          int                `json:"timeout"`            // 超时时间（秒）
	VerifySSL        bool               `json:"verify_ssl"`         // 是否验证 SSL
	Logger           *slog.Logger       `json:"-"`                  // 日志
	Proxies          RegionEnvUrls      `json:"proxies"`            // 代理
	ProxyURL         string             `json:"proxy_url"`          // HTTP 请求代理地址，为空时直连
	StaticFileServer string             `json:"static_file_server"` // 静态文件服务器地址（用于访问物流面单文件）
}
