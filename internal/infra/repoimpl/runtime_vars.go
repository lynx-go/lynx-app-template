package repoimpl

import (
	"context"
	"sync"
	"time"

	"github.com/lynx-go/lynx-app-template/internal/domain/runtimevars"
	"github.com/lynx-go/lynx-app-template/internal/infra/clients"
	entgen "github.com/lynx-go/lynx-app-template/internal/infra/ent/gen"
	"github.com/lynx-go/lynx-app-template/internal/infra/ent/gen/runtimevar"
	"github.com/lynx-go/lynx-app-template/pkg/bigid"
	"github.com/lynx-go/lynx-app-template/pkg/session"
	"github.com/lynx-go/x/encoding/json"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

type varItem struct {
	key       string
	value     string
	updatedAt time.Time
}

type runtimeVars struct {
	vars  map[string]varItem
	mu    sync.RWMutex
	data  *clients.DataClients
	idgen bigid.IDGen
}

func (r *runtimeVars) update(ctx context.Context, k string, v string, updatedAt time.Time) error {
	userId := session.CurrentUser(ctx)
	id := r.idgen.MustNextID()
	err := r.data.DB.RuntimeVar(ctx).Create().
		SetID(id.Int64()).
		SetKey(k).
		SetValue(v).
		SetCreatedAt(updatedAt).
		SetUpdatedAt(updatedAt).
		SetCreatedBy(userId.Int64()).
		SetUpdatedBy(userId.Int64()).
		OnConflict().
		Update(func(upsert *entgen.RuntimeVarUpsert) {
			upsert.SetValue(cast.ToString(v)).
				SetUpdatedAt(updatedAt).
				SetUpdatedBy(userId.Int64())
		}).
		Exec(ctx)
	return err
}

func (r *runtimeVars) getFromDb(ctx context.Context, k string) (string, time.Time, error) {
	v, err := r.data.DB.RuntimeVar(ctx).Query().Where(runtimevar.Key(k)).Only(ctx)
	if entgen.IsNotFound(err) {
		return "", time.Time{}, nil
	}
	if err != nil {
		return "", time.Time{}, err
	}
	return v.Value, v.UpdatedAt, nil
}

func (r *runtimeVars) Set(ctx context.Context, k string, v interface{}) error {
	now := time.Now()
	vs, err := cast.ToStringE(v)
	if err != nil {
		vs, err = json.MarshalToString(v)
		if err != nil {
			return err
		}
	}
	if err := r.update(ctx, k, vs, now); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.vars[k] = varItem{
		key:       k,
		value:     vs,
		updatedAt: now,
	}

	return nil
}

func (r *runtimeVars) Get(ctx context.Context, key string) (string, error) {
	r.mu.RLock()
	item, ok := r.vars[key]
	r.mu.RUnlock()
	if !ok {
		v, updatedAt, err := r.getFromDb(ctx, key)
		if err != nil {
			return "", err
		}
		if v == "" {
			return "", nil
		}
		item = varItem{
			key:       key,
			value:     v,
			updatedAt: updatedAt,
		}
		r.mu.Lock()
		r.vars[key] = item
		r.mu.Unlock()
	}
	return item.value, nil
}

func (r *runtimeVars) GetInt(ctx context.Context, key string) (int, error) {
	v, err := r.Get(ctx, key)
	if err != nil {
		return 0, err
	}
	return cast.ToIntE(v)
}

func (r *runtimeVars) GetInt32(ctx context.Context, key string) (int32, error) {
	v, err := r.Get(ctx, key)
	if err != nil {
		return 0, err
	}
	return cast.ToInt32E(v)
}

func (r *runtimeVars) GetInt64(ctx context.Context, key string) (int64, error) {
	v, err := r.Get(ctx, key)
	if err != nil {
		return 0, err
	}
	return cast.ToInt64E(v)
}

func (r *runtimeVars) GetFloat64(ctx context.Context, key string) (float64, error) {
	v, err := r.Get(ctx, key)
	if err != nil {
		return 0, err
	}
	return cast.ToFloat64E(v)
}

func (r *runtimeVars) GetFloat32(ctx context.Context, key string) (float32, error) {
	v, err := r.Get(ctx, key)
	if err != nil {
		return 0, err
	}
	return cast.ToFloat32E(v)
}

func (r *runtimeVars) GetBool(ctx context.Context, key string) (bool, error) {
	v, err := r.Get(ctx, key)
	if err != nil {
		return false, err
	}
	return cast.ToBoolE(v)
}

func (r *runtimeVars) GetAs(ctx context.Context, key string, out any) error {
	v, err := r.Get(ctx, key)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(v), out)
}

func (r *runtimeVars) Refresh(ctx context.Context) error {
	keys := lo.Keys(r.vars)

	vals, err := r.data.DB.RuntimeVar(ctx).Query().Where(runtimevar.KeyIn(keys...)).Select(runtimevar.FieldKey, runtimevar.FieldValue, runtimevar.FieldUpdatedAt).All(ctx)
	if entgen.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, v := range vals {
		r.vars[v.Key] = varItem{
			key:       v.Key,
			value:     v.Value,
			updatedAt: v.UpdatedAt,
		}
	}
	return nil
}

func NewRuntimeVars(
	data *clients.DataClients,
) runtimevars.RuntimeVars {
	return &runtimeVars{
		vars:  make(map[string]varItem),
		data:  data,
		mu:    sync.RWMutex{},
		idgen: bigid.NewIDGen(),
	}
}
