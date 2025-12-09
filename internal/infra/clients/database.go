package clients

import (
	"context"
	"time"

	entgen "github.com/lynx-go/lynx-app-template/internal/infra/ent/gen"
	configpb "github.com/lynx-go/lynx-app-template/internal/pkg/config"
	"github.com/redis/go-redis/v9"
)

type Database struct {
	*entgen.Database
}

type DataClients struct {
	DB  *Database
	RDB *redis.Client
}

func NewDatabase(driver, dialect, source string, debug bool, pool *configpb.Database_Pool) (*Database, func(), error) {
	opts := entgen.DefaultOptions()
	opts.Debug = debug
	opts.Dialect = dialect
	if pool != nil {
		if pool.MaxOpenConns > 0 {
			opts.MaxOpenConns = int(pool.MaxOpenConns)
		}
		if pool.MaxIdleConns > 0 {
			opts.MaxIdleConns = int(pool.MaxIdleConns)
		}
		if pool.ConnMaxLifetime != "" {
			if v, err := time.ParseDuration(pool.ConnMaxLifetime); err == nil {
				opts.ConnMaxLifetime = v
			}
		}
		if pool.ConnMaxIdleTime != "" {
			if v, err := time.ParseDuration(pool.ConnMaxIdleTime); err == nil {
				opts.ConnMaxIdleTime = v
			}
		}
	}

	db, closeFn, err := entgen.NewDatabase(driver, source, opts)
	if err != nil {
		return nil, nil, err
	}
	return &Database{db}, closeFn, nil
}

func NewDataClients(cfg *configpb.AppConfig) (*DataClients, func(), error) {
	c := cfg.GetData()
	ctx := context.Background()
	db, closeDb, err := NewDatabase(c.Database.Driver, c.Database.Dialect, c.Database.Source, c.Database.Debug, c.Database.GetPool())
	if err != nil {
		return nil, nil, err
	}

	rdb := newRedis(c)
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, nil, err
	}
	return &DataClients{
			DB:  db,
			RDB: rdb,
		}, func() {
			closeDb()
			_ = rdb.Close()
		}, nil
}
func newRedis(cfg *configpb.Data) *redis.Client {
	c := cfg.Redis
	return redis.NewClient(&redis.Options{
		Addr:     c.Addr,
		Password: c.Password,
		DB:       int(c.Db),
	})

}
