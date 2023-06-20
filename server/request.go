package server

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 暂不使用

// Request 封装 request 此版本内部使用gin.context
type Request interface {
	Header(key string) string
	GetHeaders() http.Header
	GetContext() context.Context
	Get(key string) (value interface{}, exists bool)
	Set(key string, value interface{})
	Copy() Request
}

// context key
type severKey struct{}

// CommonRequest 封装gin.context 为以后增加其他入口预留空间
type CommonRequest struct {
	ctx *gin.Context
}

func newCommonRequest(ginCtx *gin.Context) *CommonRequest {
	req := &CommonRequest{ctx: ginCtx}
	// 处理国际化和翻译
	return req
}

// NewContext 通过gin.Context生成一个带通用request的上下文
func NewContext(ctx context.Context, ginCtx *gin.Context) context.Context {
	return context.WithValue(ctx, severKey{}, newCommonRequest(ginCtx))
}

// FromContext 从上下文中获取通用request
func FromContext(ctx context.Context) Request {
	if ctx == nil {
		return nil
	}
	// 防止错误传递gin.context
	if _, ok := ctx.(*gin.Context); ok {
		return nil
	}
	if ctx.Value(severKey{}) == nil {
		return nil
	}

	if s, ok := ctx.Value(severKey{}).(Request); ok {
		return s
	}
	return nil
}

// GinFromContext 从上下文中获取gin.context
func GinFromContext(ctx context.Context) *gin.Context {
	if ctx == nil {
		return nil
	}
	// 防止错误传递gin.context
	if _, ok := ctx.(*gin.Context); ok {
		return nil
	}
	if ctx.Value(severKey{}) == nil {
		return nil
	}

	if s, ok := ctx.Value(severKey{}).(Request); ok {
		return s.GetContext().(*gin.Context)
	}
	return nil
}

func (r *CommonRequest) Get(key string) (value interface{}, exists bool) {
	return r.ctx.Get(key)
}

func (r *CommonRequest) Set(key string, value interface{}) {
	r.ctx.Set(key, value)
}

func (r *CommonRequest) GetHeaders() http.Header {
	if r.ctx != nil && r.ctx.Request != nil && r.ctx.Request.Header != nil {
		return r.ctx.Request.Header.Clone()
	}
	return nil
}

func (r *CommonRequest) Header(key string) string {
	if r.ctx != nil && r.ctx.Request != nil && r.ctx.Request.Header != nil {
		return r.ctx.GetHeader(key)
	}
	return ""
}

func (r *CommonRequest) GetContext() context.Context {
	return r.ctx
}

func (r *CommonRequest) Copy() Request {
	return newCommonRequest(r.ctx.Copy())
}
