package client

import (
	"context"
	"math"
	"time"
)

// BackoffFunc 重试补偿 方法
type BackoffFunc func(ctx context.Context, req *Request, attempts int) (time.Duration, error)

func exponentialBackoff(ctx context.Context, req *Request, attempts int) (time.Duration, error) {
	return do(attempts), nil
}

func do(attempts int) time.Duration {
	if attempts > 13 {
		return 2 * time.Minute
	}
	return time.Duration(math.Pow(float64(attempts), math.E)) * time.Millisecond * 100
}
