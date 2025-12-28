package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/sannticloud/sannti-cli/internal/client"
	"github.com/sannticloud/sannti-cli/internal/config"
	"github.com/sannticloud/sannti-cli/internal/models"
	"github.com/sannticloud/sannti-cli/internal/output"
)

// networkCmd represents the network command
var networkCmd = &cobra.Command{
	Use:   "network",
	Short: "Manage networks",
	Long:  `List and manage Sannti Cloud networks.`,
}

// networkListCmd lists all networks
var networkListCmd = &cobra.Command{
	Use:   "list",
	Short: "List networks",
	Long:  `List all networks, optionally filtered by region.`,
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

		networks, err := c.ListNetworks(region)
		if err != nil {
			return fmt.Errorf("failed to list networks: %w", err)
		}

		if len(networks) == 0 {
			output.PrintInfo("No networks found")
			return nil
		}

		// Convert to interface slice
		dataSlice := make([]interface{}, len(networks))
		for i, net := range networks {
			dataSlice[i] = net
		}

		return output.Print(
			dataSlice,
			output.Format(outputFormat),
			[]string{"UUID", "NAME", "STATE", "REGION", "CIDR"},
			func(item interface{}) []string {
				net := item.(models.Network)
				return []string{net.UUID, net.Name, net.State, net.ZoneName, net.Cidr}
			},
		)
	},
}

func init() {
	rootCmd.AddCommand(networkCmd)
	networkCmd.AddCommand(networkListCmd)
}
