package plog

import (
	"context"
	"github.com/yulecd/pp-common/server"
	"github.com/gin-gonic/gin"
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	InitWithPath("./", "prod")
	ginCtx := &gin.Context{}
	ctx := server.NewContext(context.Background(), ginCtx)
	Infof(ctx, "hello %s", "1111")
	Infof(nil, "hello %s", "2222")
	Infof(nil, "hello %s", "3333")
	for i := 0; i < 10; i++ {
		Infof(ctx, "hello %d", i)
		time.Sleep(time.Second)
	}
}

func TestNoInit(t *testing.T) {
	Infof(nil, "hello %s", "2222")
}

func BenchmarkLog(b *testing.B) {
	InitWithPath("./", "prod")
	ctx := &gin.Context{}
	for i := 0; i < b.N; i++ {
		Info(ctx, i)
	}
}
