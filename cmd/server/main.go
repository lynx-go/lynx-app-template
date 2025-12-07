package main

import (
	"context"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lynx-go/lynx"
	"github.com/lynx-go/lynx/contrib/zap"
)

var (
	version string
)

func main() {
	opts := lynx.NewOptions(
		lynx.WithName("artizard-api"),
		lynx.WithVersion(version),
		lynx.WithUseDefaultConfigFlagsFunc(),
	)
	app := lynx.New(opts, func(ctx context.Context, app lynx.Lynx) error {
		app.SetLogger(zap.MustNewLogger(app))

		boot, cleanup, err := wireBootstrap(app, app.Logger())
		if err != nil {
			log.Fatal(err)
		}
		if err := app.Hook(lynx.OnStop(func(ctx context.Context) error {
			cleanup()
			return nil
		})); err != nil {
			return err
		}
		return boot.Build(app)
	})
	app.Run()
}
