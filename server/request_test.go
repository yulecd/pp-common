package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"testing"
)

func TestContext(t *testing.T) {
	c := context.Background()
	ginC := &gin.Context{}
	ginC.Set("a", 1)
	ctx := NewContext(c, ginC)
	getGinC := GinFromContext(ctx)
	if getGinC == nil {
		t.Error("get gin context fail")
	}
	r, ok := getGinC.Get("a")
	if !ok {
		t.Error("get gin context value fail")
	}
	rr := r.(int)
	if rr != 1 {
		t.Error("get gin context value wrong")
	}
}
