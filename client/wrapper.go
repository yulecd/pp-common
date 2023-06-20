package client

import (
	"context"
)

// WrapperChain client 中间件
type WrapperChain func(next Wrapper) Wrapper

type Wrapper func(ctx context.Context, req *Request) (*Response, error)
