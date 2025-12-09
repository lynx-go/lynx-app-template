package cmd

import (
	"context"
	"time"

	"github.com/lynx-go/lynx"
	"github.com/lynx-go/lynx-app-template/internal/infra/server"
	configpb "github.com/lynx-go/lynx-app-template/internal/pkg/config"
	"github.com/lynx-go/lynx-app-template/pkg/pubsub"
	"github.com/lynx-go/lynx/contrib/zap"
	"github.com/lynx-go/x/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type CmdArgs struct {
	Cmd  *cobra.Command
	Args []string
}

type AppContext struct {
	App    lynx.Lynx
	PubSub *pubsub.PubSub
}

func buildCli(cmd *cobra.Command, args []string, fn func(ctx context.Context, appCtx *AppContext, cmdArgs *CmdArgs) error, sleepDuration time.Duration) *lynx.CLI {
	return lynx.New(newOptionsFromCmd(cmd), func(ctx context.Context, app lynx.Lynx) error {
		app.SetLogger(zap.MustNewLogger(app))
		config := &configpb.AppConfig{}
		if err := app.Config().Unmarshal(config, lynx.TagNameJSON); err != nil {
			return err
		}
		pubSub := server.NewPubSub()
		// CLI 禁用 kafka 监听
		binder := server.NewKafkaBinderForCli(pubSub, config)
		if err := app.Hook(lynx.Components(pubSub, binder), lynx.ComponentBuilders(binder.Builders()...)); err != nil {
			return err
		}

		return app.CLI(func(ctx context.Context) error {
			if err := fn(ctx, &AppContext{App: app, PubSub: pubSub}, &CmdArgs{Cmd: cmd, Args: args}); err != nil {
				return err
			}
			if sleepDuration > 0 {
				// wait pubsub completed
				log.InfoContext(ctx, "waiting 1 seconds for pubsub completed")
				time.Sleep(sleepDuration)
			}
			return nil
		})
	})
}
func newOptionsFromCmd(cmd *cobra.Command) *lynx.Options {
	return lynx.NewOptions(
		lynx.WithName(cmd.Root().Name()+":"+cmd.Name()),
		lynx.WithBindConfigFunc(func(f *pflag.FlagSet, v *viper.Viper) error {
			if cd, _ := cmd.Root().PersistentFlags().GetString("config-dir"); cd != "" {
				v.AddConfigPath(cd)
			}

			return nil
		}),
	)
}
