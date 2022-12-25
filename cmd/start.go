package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stackup-wallet/stackup-bundler/internal/start"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts an instance in the specified mode",
	Long: `The start command has the following modes:
	
	1. private: A JSON-RPC client and bundler backed by a private mempool.`,
	Run: func(cmd *cobra.Command, args []string) {
		if viper.GetString("mode") == "private" {
			start.PrivateMode()
		} else {
			panic(fmt.Sprintf("Fatal flag error: \"%s\" mode not supported", viper.GetString("mode")))
		}
	},
}

var mode string

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringVarP(&mode, "mode", "m", "", "Required. See acceptable values above.")
	if err := startCmd.MarkFlagRequired("mode"); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("mode", startCmd.Flags().Lookup("mode")); err != nil {
		panic(err)
	}
}
