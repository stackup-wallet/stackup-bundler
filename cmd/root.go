package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "stackup-bundler",
	Short: "Standalone ERC-4337 Client & Bundler",
	Long:  "A standalone RPC client and bundler for relaying UserOperations to the EntryPoint.",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {}
