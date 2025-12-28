package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/sannticloud/sannti-cli/internal/client"
	"github.com/sannticloud/sannti-cli/internal/config"
	"github.com/sannticloud/sannti-cli/internal/models"
	"github.com/sannticloud/sannti-cli/internal/output"
)

// regionCmd represents the region command
var regionCmd = &cobra.Command{
	Use:   "region",
	Short: "Manage regions",
	Long:  `List and manage Sannti Cloud regions.`,
}

// regionListCmd represents the region list command
var regionListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available regions",
	Long:  `Display all available Sannti Cloud regions with their details.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load config
		cfg, err := config.LoadConfig()
		if err != nil {
			return err
		}

		// Create client
		c := client.NewClient(cfg.AccessKey, cfg.SecretKey)

		// List zones (regions)
		zones, err := c.ListZones()
		if err != nil {
			return fmt.Errorf("failed to list regions: %w", err)
		}

		// Filter only active zones
		var activeZones []models.Zone
		for _, zone := range zones {
			if zone.IsActive {
				activeZones = append(activeZones, zone)
			}
		}

		if len(activeZones) == 0 {
			output.PrintInfo("No active regions found")
			return nil
		}

		// Convert to interface slice for output
		dataSlice := make([]interface{}, len(activeZones))
		for i, z := range activeZones {
			dataSlice[i] = z
		}

		// Print output
		return output.Print(
			dataSlice,
			output.Format(outputFormat),
			[]string{"REGION", "COUNTRY", "STATUS"},
			func(item interface{}) []string {
				z := item.(models.Zone)
				status := "inactive"
				if z.IsActive {
					status = "active"
				}
				return []string{z.Name, z.CountryName, status}
			},
		)
	},
}

func init() {
	rootCmd.AddCommand(regionCmd)
	regionCmd.AddCommand(regionListCmd)
}
