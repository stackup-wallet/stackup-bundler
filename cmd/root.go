package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "stackup-bundler",
	Short: "ERC-4337 Bundler",
	Long:  "A modular Go implementation of an ERC-4337 Bundler.",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {}
