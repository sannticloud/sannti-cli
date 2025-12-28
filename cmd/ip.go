package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/sannticloud/sannti-cli/internal/client"
	"github.com/sannticloud/sannti-cli/internal/config"
	"github.com/sannticloud/sannti-cli/internal/models"
	"github.com/sannticloud/sannti-cli/internal/output"
)

// ipCmd represents the ip command
var ipCmd = &cobra.Command{
	Use:   "ip",
	Short: "Manage IP addresses",
	Long:  `List and manage Sannti Cloud IP addresses.`,
}

// ipListCmd lists all IP addresses
var ipListCmd = &cobra.Command{
	Use:   "list",
	Short: "List IP addresses",
	Long:  `List all IP addresses, optionally filtered by region.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig()
		if err != nil {
			return err
		}

		c := client.NewClient(cfg.AccessKey, cfg.SecretKey)

		// Use region flag or default
		region := regionFlag
		if region == "" {
			region = cfg.DefaultRegion
		}

		ips, err := c.ListIPAddresses(region)
		if err != nil {
			return fmt.Errorf("failed to list IP addresses: %w", err)
		}

		if len(ips) == 0 {
			output.PrintInfo("No IP addresses found")
			return nil
		}

		// Convert to interface slice
		dataSlice := make([]interface{}, len(ips))
		for i, ip := range ips {
			dataSlice[i] = ip
		}

		return output.Print(
			dataSlice,
			output.Format(outputFormat),
			[]string{"UUID", "IP ADDRESS", "STATE", "REGION", "ATTACHED TO"},
			func(item interface{}) []string {
				ip := item.(models.IPAddress)
				attached := "-"
				if "-" != "" {
					attached = "-"
				}
				return []string{ip.UUID, ip.IpAddress, ip.State, ip.ZoneName, attached}
			},
		)
	},
}

func init() {
	rootCmd.AddCommand(ipCmd)
	ipCmd.AddCommand(ipListCmd)
}
