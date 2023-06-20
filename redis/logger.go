package redis

import (
	"context"
	"time"

	"github.com/yulecd/pp-common/plog"

	"github.com/go-redis/redis/v8"
)

type RedisLogger struct {
}

func (t RedisLogger) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	if cmd.Name() == `ping` {
		return ctx, nil
	}
	return context.WithValue(ctx, `log_cache_Process_start`, time.Now()), nil
}

func (t RedisLogger) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	if cmd.Name() == `ping` {
		return nil
	}
	s, ok := ctx.Value(`log_cache_Process_start`).(time.Time)
	if ok {
		plog.Infof(ctx, "redis cmd: %s const:%v", cmd.Name(), s)
	}
	return nil
}

func (t RedisLogger) BeforeProcessPipeline(ctx context.Context, _ []redis.Cmder) (context.Context, error) {
	return context.WithValue(ctx, `log_cache_ProcessPipeline_start`, time.Now()), nil
}

func (t RedisLogger) AfterProcessPipeline(ctx context.Context, cmdList []redis.Cmder) error {
	s, ok := ctx.Value(`log_cache_ProcessPipeline_start`).(time.Time)
	if ok {
		plog.Infof(ctx, "redis cmd: %v const:%v", cmdList, s)
	}
	return nil
}
