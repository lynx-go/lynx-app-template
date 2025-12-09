/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"

	"github.com/lynx-go/lynx"
	"github.com/lynx-go/x/encoding/json"
	"github.com/lynx-go/x/log"
	"github.com/spf13/cobra"
)

// printConfigCmd represents the printConfig command
var printConfigCmd = &cobra.Command{
	Use:   "printConfig",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		cli := buildCli(cmd, args, func(ctx context.Context, appCtx *AppContext, cmdArgs *CmdArgs) error {
			config := map[string]any{}
			if err := appCtx.App.Config().Unmarshal(&config, lynx.TagNameJSON); err != nil {
				return err
			}
			log.InfoContext(ctx, "print config", "configs", json.MustMarshalToString(config))
			return nil
		}, 0)
		cli.Run()
	},
}

func init() {
	rootCmd.AddCommand(printConfigCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// printConfigCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// printConfigCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
