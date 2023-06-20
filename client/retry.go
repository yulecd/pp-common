package client

import "context"

// RetryFunc 是否需要重试
type RetryFunc func(ctx context.Context, req *Request, resp *Response, retryCount int, err error) (bool, error)

func RetryOnError(ctx context.Context, req *Request, resp *Response, retryCount int, err error) (bool, error) {
	if err == nil {
		return false, nil
	}

	if resp.hresp == nil {
		return false, err
	}

	code := resp.GetStatusCode()

	switch code {
	case 408, 500:
		return true, nil
	default:
		return false, nil
	}
}
