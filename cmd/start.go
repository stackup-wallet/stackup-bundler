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
	
	1. private: A bundler backed by a private mempool and compatible with all EVM networks.
	2. searcher: A bundler backed by the P2P mempool and integrated with a Block Builder API.`,
	Run: func(cmd *cobra.Command, args []string) {
		if viper.GetString("mode") == "private" {
			start.PrivateMode()
		} else if viper.GetString("mode") == "searcher" {
			start.SearcherMode()
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
