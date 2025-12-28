package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/sannticloud/sannti-cli/internal/client"
	"github.com/sannticloud/sannti-cli/internal/config"
	"github.com/sannticloud/sannti-cli/internal/models"
	"github.com/sannticloud/sannti-cli/internal/output"
)

// firewallCmd represents the firewall command
var firewallCmd = &cobra.Command{
	Use:   "firewall",
	Short: "Manage firewall rules",
	Long:  `List and manage Sannti Cloud firewall rules.`,
}

// firewallListCmd lists all firewall rules
var firewallListCmd = &cobra.Command{
	Use:   "list",
	Short: "List firewall rules",
	Long:  `List all firewall rules.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig()
		if err != nil {
			return err
		}

		c := client.NewClient(cfg.AccessKey, cfg.SecretKey)

		region := regionFlag
		if region == "" {
			region = cfg.DefaultRegion
		}

		rules, err := c.ListFirewallRules(region)
		if err != nil {
			return fmt.Errorf("failed to list firewall rules: %w", err)
		}

		if len(rules) == 0 {
			output.PrintInfo("No firewall rules found")
			return nil
		}

		// Convert to interface slice
		dataSlice := make([]interface{}, len(rules))
		for i, rule := range rules {
			dataSlice[i] = rule
		}

		return output.Print(
			dataSlice,
			output.Format(outputFormat),
			[]string{"UUID", "PROTOCOL", "PORT RANGE", "CIDR", "STATE"},
			func(item interface{}) []string {
				rule := item.(models.FirewallRule)
				portRange := fmt.Sprintf("%s-%s", rule.StartPort, rule.EndPort)
				if rule.StartPort == rule.EndPort {
					portRange = rule.StartPort
				}
				return []string{rule.UUID, rule.Protocol, portRange, rule.CidrList, rule.State}
			},
		)
	},
}

func init() {
	rootCmd.AddCommand(firewallCmd)
	firewallCmd.AddCommand(firewallListCmd)
}
