package client

import (
	"time"
)

type Options struct {
	debug      bool
	timeout    time.Duration
	backoff    BackoffFunc
	retry      RetryFunc
	retries    int
	wrappers   []WrapperChain
	BaseURI    string
	Query      interface{}
	Headers    map[string]interface{}
	FormParams map[string]interface{}
	JSON       interface{}
}

// NewOptions instances default Options
func NewOptions(options ...Option) Options {
	opts := Options{
		debug:    false,
		timeout:  DefaultTimeout,
		backoff:  DefaultBackoff,
		retry:    DefaultRetry,
		retries:  DefaultRetries,
		wrappers: DefaultWrappers,
	}

	for _, o := range options {
		o(&opts)
	}

	return opts
}

func (o *Options) merge(opts Options) {
	if opts.BaseURI != "" {
		o.BaseURI = opts.BaseURI
	}
	if opts.Query != nil {
		o.Query = opts.Query
	}
	if opts.FormParams != nil {
		o.FormParams = opts.FormParams
	}
	if opts.JSON != nil {
		o.JSON = opts.JSON
	}
	if opts.Headers != nil {
		o.Headers = opts.Headers
	}
}

// Debug logs request vebose info
func Debug(d bool) Option {
	return func(o *Options) {
		o.debug = d
	}
}

// The request base uri.
func BaseURI(uri string) Option {
	return func(o *Options) {
		o.BaseURI = uri
	}
}

// Adds a WrapperChain to a list of options passed into the client
func Wrap(w WrapperChain) Option {
	return func(o *Options) {
		o.wrappers = append(o.wrappers, w)
	}
}

// The request timeout.
func Timeout(t time.Duration) Option {
	return func(o *Options) {
		o.timeout = t
	}
}

// Retry sets the retry function to be used when re-trying.
func Retry(fn RetryFunc) Option {
	return func(o *Options) {
		o.retry = fn
	}
}

// Number of retries when making the request.
func Retries(i int) Option {
	return func(o *Options) {
		o.retries = i
	}
}
