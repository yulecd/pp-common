package client

import (
	"fmt"
	"time"

	"github.com/yulecd/pp-common/config"
)

const (
	ServiceConfigName = "service"
)

type Config struct {
	Name    string            `json:"name" yaml:"name"`
	BaseUri string            `json:"base_uri" yaml:"base_uri"`
	Timeout int               `json:"timeout" yaml:"timeout"` // 超时时间 单位毫秒
	Headers map[string]string `json:"headers" yaml:"headers"` // header 头
}

// NewHttpClientWithConfig 根据config配置初始化client
func NewHttpClientWithConfig(name string, service string, opt ...Option) *Request {
	nOpt := make([]Option, 0, len(opt)+5)

	var cm map[string]Config
	err := config.Load(name, &cm)
	if err == nil {
		if c, ok := cm[service]; ok {
			nOpt = append(nOpt, func(options *Options) {
				if len(c.BaseUri) > 0 {
					options.BaseURI = c.BaseUri
				}
				if c.Timeout > 0 {
					options.timeout = time.Duration(c.Timeout) * time.Microsecond
				}
				if len(c.Headers) > 0 {
					if options.Headers == nil {
						options.Headers = make(map[string]interface{})
					}
					for key, value := range c.Headers {
						options.Headers[key] = value
					}
				}
			})
		}
	} else {
		fmt.Println("newHttpClientWithConfig error:", err)
	}
	nOpt = append(nOpt, opt...)
	return newHttpClient(nOpt...)
}
