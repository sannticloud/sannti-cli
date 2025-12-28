package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/sannticloud/sannti-cli/internal/config"
	"github.com/sannticloud/sannti-cli/internal/output"
)

var (
	outputFormat string
	regionFlag   string
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "sannti",
	Short: "Sannti Cloud CLI - Beyond Cloud. Without Barriers.",
	Long: `Sannti CLI is a command-line interface for managing Sannti Cloud resources.

The CLI provides access to compute instances, Kubernetes clusters, networking,
and other cloud resources through an intuitive command structure.

To get started, run: sannti configure`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		output.PrintError(err.Error())
		os.Exit(1)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "table", "Output format (table, json, yaml)")
	rootCmd.PersistentFlags().StringVarP(&regionFlag, "region", "r", "", "Region (overrides default)")

	// Initialize config
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	if err := config.InitConfig(); err != nil {
		// Only show error if not running configure command
		if rootCmd.CalledAs() != "configure" && rootCmd.CalledAs() != "version" {
			fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
		}
	}
}
