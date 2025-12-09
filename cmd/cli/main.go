package main

import (
	"context"

	"github.com/lynx-go/lynx"
	"github.com/lynx-go/lynx/contrib/zap"
	"github.com/lynx-go/x/log"
)

var (
	version string
)

func main() {

	opts := lynx.NewOptions(
		lynx.WithName("lynx-app-cli"),
		lynx.WithVersion(version),
		lynx.WithUseDefaultConfigFlagsFunc(),
	)
	cli := lynx.New(opts, func(ctx context.Context, app lynx.Lynx) error {
		app.SetLogger(zap.MustNewLogger(app))

		return app.CLI(func(ctx context.Context) error {
			log.InfoContext(ctx, "hello lynx cli")
			return nil
		})
	})
	cli.Run()
}
