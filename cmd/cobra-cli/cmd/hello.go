/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"time"

	"github.com/lynx-go/lynx-app-template/internal/domain/events"
	"github.com/lynx-go/x/log"
	"github.com/spf13/cobra"
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
		cli := buildCli(cmd, args, func(ctx context.Context, appCtx *AppContext, cmdArgs *CmdArgs) error {
			toUid, _ := cmdArgs.Cmd.Flags().GetInt("to")
			if err := appCtx.PubSub.Publish(ctx, events.ProducerNameHello, "hello", &events.HelloEvent{
				Message: "Hello " + args[0],
				Time:    time.Now(),
			}); err != nil {
				log.ErrorContext(ctx, "publish hello event error", err)
			}
			log.InfoContext(ctx, "hello lynx cli", "args", args, "to_uid", toUid)
			//time.Sleep(1 * time.Second)
			return nil
		}, 1*time.Second)
		cli.Run()
	},
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
