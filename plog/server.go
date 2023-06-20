package plog

// 通用请求 框架中使用log

import (
	"context"
	"github.com/yulecd/pp-common/server"
	"github.com/yulecd/pp-common/trace"
	"github.com/gin-gonic/gin"
	"time"
)

const (
	ContextLogKey = "__context_log_key__"
)

const (
	requestStartTimeKey = "start"
	traceIDKey          = "trace"
	queryPathKey        = "path"
	nameKey             = "name"
)

// GetDefaultFieldEntry 获取带默认字段的日志入口
// 依次从 ctx header 获取traceId
func GetDefaultFieldEntry(ctx context.Context) *Entry {
	var logEntry *Entry
	var isNewTrace bool
	// 从上下文获取log对象
	request := server.FromContext(ctx)
	if request != nil {
		logEntryValue, ok := request.Get(ContextLogKey)
		if ok && logEntryValue != nil {
			logEntry, ok = logEntryValue.(*Entry)
			if ok {
				return logEntry
			}
		}
	}

	// 初始化
	traceID := ""
	if request != nil {
		// 尝试从上下文获取
		traceIDInterface, ok := request.Get(trace.ContextTraceId)
		if ok {
			if traceIDStr, ok := traceIDInterface.(string); ok {
				traceID = traceIDStr
			}
		}
	}
	if traceID == "" && request != nil {
		// 尝试从header获取
		traceID = request.Header(trace.HeaderTraceIdKey)
		if len(traceID) > 50 {
			traceID = traceID[:50]
		}
	}
	if traceID == "" {
		isNewTrace = true
		traceID = trace.ID()
	}
	// 此处依赖了gin.Context
	ginCtx := server.GinFromContext(ctx)
	path := ""
	if ginCtx != nil && ginCtx.Request != nil && ginCtx.Request.URL != nil {
		path = ginCtx.Request.URL.Path
	}

	// 默认字段
	fields := map[string]interface{}{
		requestStartTimeKey: time.Now().Format(logTimeFormatter),
		traceIDKey:          traceID,
		queryPathKey:        path,
	}
	logEntry = stdLogger.withFields(fields)

	// 保存log对象到上下文
	if request != nil {
		request.Set(ContextLogKey, logEntry)
	}

	// 保存traceID
	if isNewTrace {
		if request != nil {
			request.Set(trace.ContextTraceId, traceID)
		}
		// 此处依赖了gin.Context
		if ginCtx != nil && ginCtx.Request != nil && ginCtx.Request.Header != nil {
			ginCtx.Header(trace.HeaderTraceIdKey, traceID)
		}
	}

	return logEntry
}

// FlushLogEntry 清除log对象
func FlushLogEntry(ctx context.Context) {
	request := server.FromContext(ctx)
	if request != nil {
		request.Set(ContextLogKey, nil)
	}
}

// GetDefaultFieldEntryFromGin get log from gin.context
func GetDefaultFieldEntryFromGin(c *gin.Context) *Entry {
	return GetDefaultFieldEntry(server.NewContext(context.Background(), c))
}

// getFromRequest return a Entry, you should use Entry in api handler to log
func getFromRequest(c context.Context) *Entry {
	var (
		e  *Entry
		ok bool
	)
	if c == nil {
		return nil
	}

	request := server.FromContext(c)
	if request == nil {
		return nil
	}
	ee, ok := request.Get(ContextLogKey)
	if ok && ee != nil {
		e, ok = ee.(*Entry)
		if ok {
			return e
		}
	}
	return nil
}

// warp stdLogger
func Debug(ctx context.Context, args ...interface{}) {
	GetDefaultFieldEntry(ctx).Debug(args...)
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	GetDefaultFieldEntry(ctx).Debugf(format, args...)
}

func Info(ctx context.Context, args ...interface{}) {
	GetDefaultFieldEntry(ctx).Info(args...)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	GetDefaultFieldEntry(ctx).Infof(format, args...)
}

func Warn(ctx context.Context, args ...interface{}) {
	GetDefaultFieldEntry(ctx).Warn(args...)
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	GetDefaultFieldEntry(ctx).Warnf(format, args...)
}

func Error(ctx context.Context, args ...interface{}) {
	GetDefaultFieldEntry(ctx).Error(args...)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	GetDefaultFieldEntry(ctx).Errorf(format, args...)
}

// Fatal will call os.Exit(1), be careful to use.
func Fatal(ctx context.Context, args ...interface{}) {
	GetDefaultFieldEntry(ctx).Fatal(args...)
}

// Fatalf will call os.Exit(1), be careful to use.
func Fatalf(ctx context.Context, format string, args ...interface{}) {
	GetDefaultFieldEntry(ctx).Fatalf(format, args...)
}
