package config

import (
	"log/slog"
	"time"
)

type URLPair struct {
	Prod string
	Test string
}

type Config struct {
	Env              string             `json:"env"`                // 环境（dev：开发环境、test：测试环境、prod：生产环境）
	Debug            bool               `json:"debug"`              // 是否为调试模式
	Region           string             `json:"region"`             // 区域
	AppKey           string             `json:"app_key"`            // App Key
	AppSecret        string             `json:"app_secret"`         // App 秘钥
	AccessToken      string             `json:"access_token"`       // Access Token
	Timeout          time.Duration      `json:"timeout"`            // 超时时间（秒）
	VerifySSL        bool               `json:"verify_ssl"`         // 是否验证 SSL
	Logger           *slog.Logger       `json:"-"`                  // 日志
	OverwriteUrls    map[string]URLPair `json:"overwrite_urls"`     // 覆盖 URL
	StaticFileServer string             `json:"static_file_server"` // 静态文件服务器地址（用于访问物流面单文件）
}
