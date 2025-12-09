/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"

	"github.com/lynx-go/lynx"
	"github.com/lynx-go/lynx/contrib/zap"
	"github.com/lynx-go/x/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// helloCmd represents the hello command
var helloCmd = &cobra.Command{
	Use:   "hello",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cli := buildCli(cmd, args, func(ctx context.Context, app lynx.Lynx, cmd *cobra.Command, args []string) error {
			toUid, _ := cmd.Flags().GetInt("to")
			log.InfoContext(ctx, "hello lynx cli", "args", args, "to_uid", toUid)
			return nil
		})
		cli.Run()
	},
}

func buildCli(cmd *cobra.Command, args []string, fn func(ctx context.Context, app lynx.Lynx, cmd *cobra.Command, args []string) error) *lynx.CLI {
	return lynx.New(newOptionsFromCmd(cmd), func(ctx context.Context, app lynx.Lynx) error {
		app.SetLogger(zap.MustNewLogger(app))

		return app.CLI(func(ctx context.Context) error {
			return fn(ctx, app, cmd, args)
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

func init() {
	rootCmd.AddCommand(helloCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// helloCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// helloCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	helloCmd.Flags().IntP("to", "t", 0, "hello to uid")
}
