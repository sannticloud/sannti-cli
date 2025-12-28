package cmd

import (
"fmt"

"github.com/spf13/cobra"

"github.com/sannticloud/sannti-cli/internal/client"
"github.com/sannticloud/sannti-cli/internal/config"
"github.com/sannticloud/sannti-cli/internal/models"
"github.com/sannticloud/sannti-cli/internal/output"
)

// k8sCmd represents the kubernetes command
var k8sCmd = &cobra.Command{
Use:     "k8s",
Aliases: []string{"kubernetes"},
Short:   "Manage Kubernetes clusters",
Long:    `List available Kubernetes versions for cluster creation.`,
}

// k8sVersionsCmd lists available Kubernetes versions
var k8sVersionsCmd = &cobra.Command{
Use:   "versions",
Short: "List available Kubernetes versions",
Long:  `List all available Kubernetes versions that can be used to create clusters.`,
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

versions, err := c.ListKubernetesVersions(region)
if err != nil {
return fmt.Errorf("failed to list Kubernetes versions: %w", err)
}

if len(versions) == 0 {
output.PrintInfo("No Kubernetes versions found")
return nil
}

// Convert to interface slice
dataSlice := make([]interface{}, len(versions))
for i, v := range versions {
dataSlice[i] = v
}

return output.Print(
dataSlice,
output.Format(outputFormat),
[]string{"UUID", "VERSION", "MIN CPU", "MIN MEMORY (MB)"},
func(item interface{}) []string {
v := item.(models.KubernetesVersion)
return []string{
v.UUID,
v.Name,
fmt.Sprintf("%d", v.MinCPUNumber),
fmt.Sprintf("%d", v.MinMemory),
}
},
)
},
}

func init() {
rootCmd.AddCommand(k8sCmd)
k8sCmd.AddCommand(k8sVersionsCmd)
}
