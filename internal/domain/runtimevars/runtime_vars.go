package runtimevars

import "context"

type RuntimeVars interface {
	Setter
	Getter
	Scheduler
}

type Getter interface {
	Set(ctx context.Context, k string, v interface{}) error
	Get(ctx context.Context, name string) (string, error)
	GetInt(ctx context.Context, name string) (int, error)
	GetInt32(ctx context.Context, name string) (int32, error)
	GetInt64(ctx context.Context, name string) (int64, error)
	GetFloat64(ctx context.Context, name string) (float64, error)
	GetFloat32(ctx context.Context, name string) (float32, error)
	GetBool(ctx context.Context, name string) (bool, error)
	GetAs(ctx context.Context, name string, out any) error
}

type Setter interface {
	Set(ctx context.Context, k string, v interface{}) error
}

type Scheduler interface {
	Refresh(ctx context.Context) error
}
