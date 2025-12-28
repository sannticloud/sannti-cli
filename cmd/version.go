package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	Version = "v0.1.0"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show Sannti CLI version",
	Long:  `Display the current version of the Sannti CLI.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Sannti CLI %s\n", Version)
		fmt.Println("Built for Sannti Cloud - Beyond Cloud. Without Barriers.")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
