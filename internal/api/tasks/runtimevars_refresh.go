package tasks

import (
	"context"

	"github.com/lynx-go/lynx-app-template/internal/domain/runtimevars"
	"github.com/lynx-go/lynx-app-template/internal/pkg/varkeys"
	"github.com/lynx-go/lynx/contrib/schedule"
	"github.com/lynx-go/x/log"
)

func NewRuntimeVarsRefresh(
	runtimeVars runtimevars.RuntimeVars,
) *RuntimeVarsRefresh {
	task := &RuntimeVarsRefresh{
		runtimeVars: runtimeVars,
	}
	_ = task.preload(context.TODO())
	return task
}

type RuntimeVarsRefresh struct {
	runtimeVars runtimevars.RuntimeVars
}

func (r *RuntimeVarsRefresh) Name() string {
	return "RuntimeVarsRefresh"
}

func (r *RuntimeVarsRefresh) Cron() string {
	return "@every 1m"
}

func (r *RuntimeVarsRefresh) HandlerFunc() schedule.HandlerFunc {
	return func(ctx context.Context) error {
		return r.runtimeVars.Refresh(ctx)
	}
}

// preload 预加载数据
func (r *RuntimeVarsRefresh) preload(ctx context.Context) error {
	log.DebugContext(ctx, "preload runtime vars")
	_, _ = r.runtimeVars.Get(ctx, varkeys.ApiKey)
	_, _ = r.runtimeVars.Get(ctx, varkeys.AppInitialized)
	return nil
}

var _ schedule.Task = new(RuntimeVarsRefresh)
