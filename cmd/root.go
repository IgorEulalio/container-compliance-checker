package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var configPath string

var rootCmd = &cobra.Command{
	Use:   "container-compliance-checker",
	Short: "Performs compliance checks for containers",
	Long:  `A CLI tool to perform various compliance checks on containers.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&configPath,
		"config",
		"c",
		"./config.yaml",
		"Path to the config file",
	)
}
