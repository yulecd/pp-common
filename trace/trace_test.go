package trace

import (
	"context"
	"testing"

	"github.com/yulecd/pp-common/server"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestID(t *testing.T) {
	InitGenerator()
	id := ID()
	if len(id) == 0 {
		t.Errorf("ID fail")
	}
}

func BenchmarkID(b *testing.B) {
	InitGenerator()
	for i := 0; i < b.N; i++ {
		ID()
	}
}

func BenchmarkIDB(b *testing.B) {
	InitGenerator()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ID()
		}
	})
}

func TestFlushTraceId(t *testing.T) {
	ctx := server.NewContext(context.Background(), &gin.Context{})
	id := GetTraceIdFromContext(ctx)
	assert.NotNil(t, id)
	id2 := GetTraceIdFromContext(ctx)
	assert.Equal(t, id, id2)
	id3 := "12345678"
	SetTraceId(ctx, id3)
	id4 := GetTraceIdFromContext(ctx)
	assert.Equal(t, id3, id4)
}
