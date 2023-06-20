package trace

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/yulecd/pp-common/server"
)

const (
	ContextTraceId   = "__context_trace_id__"
	HeaderTraceIdKey = "x-trace-id"
)

const (
	seqPoolLen = 8000
	invalidIP  = "00000000"
)

var generator = struct {
	sync.Mutex
	seqPool []string
	ip      string
	pid     string
	i       int
}{}

func InitGenerator() {
	generator.seqPool = make([]string, seqPoolLen)
	generator.ip = GetLocalIP()
	if generator.ip == "" {
		generator.ip = invalidIP
	}
	// pid in mac can be 99998(0x1869E)
	generator.pid = fmt.Sprintf("%04x", os.Getpid())[:4]
	for i := 0; i < seqPoolLen; i++ {
		generator.seqPool[i] = fmt.Sprintf("%04x", i+1)
	}
}

func ID() string {
	generator.Lock()
	// 未初始化时 仍可使用
	seqID := "0000"
	if len(generator.seqPool) > 0 {
		seqID = generator.seqPool[generator.i]
	}
	generator.i++
	if generator.i >= seqPoolLen {
		generator.i = 0
	}
	generator.Unlock()
	timestamp := fmt.Sprintf("%016x", time.Now().UnixNano()/int64(time.Millisecond))
	return seqID + timestamp + generator.ip + generator.pid
}

// GetTraceIdFromContext 从上下文中获取traceId 不存在的话生成一个保存到上下文中
func GetTraceIdFromContext(ctx context.Context) string {
	// 获取traceId
	traceID := ""
	request := server.FromContext(ctx)
	if request != nil {
		// 上下文中获取
		obj, ok := request.Get(ContextTraceId)
		if ok {
			traceID = obj.(string)
		}
	}

	if traceID == "" {
		traceID = ID()
		// 保存回去 供后续使用
		if request != nil {
			request.Set(ContextTraceId, traceID)
		}
	}
	return traceID
}

// FlushTraceId 清除traceId
func FlushTraceId(ctx context.Context) {
	request := server.FromContext(ctx)
	if request != nil {
		request.Set(ContextTraceId, "")
	}
}

// SetTraceId 设置traceId
func SetTraceId(ctx context.Context, traceId string) {
	if len(traceId) > 50 {
		traceId = traceId[:50]
	}
	request := server.FromContext(ctx)
	if request != nil {
		request.Set(ContextTraceId, traceId)
	}
}
