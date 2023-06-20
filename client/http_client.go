package client

import "context"

// newHttpClient instances a http request
func newHttpClient(opt ...Option) *Request {
	opts := NewOptions(opt...)

	req := &Request{
		opts: opts,
	}

	req.WithContext(context.Background())

	return req
}
