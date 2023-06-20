package client

import (
	"time"
)

// Client 请求client 重点是option配置数据
type Client interface{}

type Option func(*Options)

var (
	DefaultClient   *Request = newHttpClient()
	DefaultBackoff           = exponentialBackoff
	DefaultRetry             = RetryOnError
	DefaultWrappers          = make([]WrapperChain, 0)
	DefaultRetries           = 1
	DefaultTimeout           = time.Second * 30

	NewClient func(...Option) *Request = newHttpClient
)

// WithDefaultRetries sets default retries number.
func WithDefaultRetries(i int) {
	DefaultRetries = i
}

// AddDefaultWrappers adds default global wrapper chain
func AddDefaultWrappers(chain WrapperChain) {
	DefaultWrappers = append(DefaultWrappers, chain)
}

// WithDefaultRetry sets default retry function.
func WithDefaultRetry(fn RetryFunc) {
	DefaultRetry = fn
}

// WithDefaultTimeout sets default timeout of request.
func WithDefaultTimeout(t time.Duration) {
	DefaultTimeout = t
}

// WithDefaultBackoff sets default backoff function policy.
func WithDefaultBackoff(fn BackoffFunc) {
	DefaultBackoff = fn
}
