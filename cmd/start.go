package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stackup-wallet/stackup-bundler/internal/start"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts an instance in either client or bundler mode",
	Long: `The start command has two modes:
	
	1. client: A JSON-RPC server for adding ops to the mempool.
	2. bundler: Workers that create batches from ops in the mempool.`,
	Run: func(cmd *cobra.Command, args []string) {
		if viper.GetString("mode") == "client" {
			start.ClientMode()
		} else if viper.GetString("mode") == "bundler" {
			fmt.Println("TODO: Implement bundler mode")
		} else {
			panic(fmt.Sprintf("Fatal flag error: \"%s\" mode not supported", viper.GetString("mode")))
		}
	},
}

var mode string

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringVarP(&mode, "mode", "m", "", "Accepts either 'client' or 'bundler'")
	startCmd.MarkFlagRequired("mode")
	viper.BindPFlag("mode", startCmd.Flags().Lookup("mode"))
}
