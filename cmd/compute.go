package cmd

import (
"fmt"

"github.com/spf13/cobra"
"github.com/sannticloud/sannti-cli/internal/client"
"github.com/sannticloud/sannti-cli/internal/config"
"github.com/sannticloud/sannti-cli/internal/models"
"github.com/sannticloud/sannti-cli/internal/output"
)

// computeCmd represents the compute command
var computeCmd = &cobra.Command{
Use:   "compute",
Short: "Manage compute instances",
Long:  `Create, list, and manage Sannti Cloud compute instances (virtual machines).`,
}

// computeListCmd lists all instances
var computeListCmd = &cobra.Command{
Use:   "list",
Short: "List compute instances",
Long:  `List all compute instances, optionally filtered by region.`,
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

instances, err := c.ListInstances(region)
if err != nil {
return fmt.Errorf("failed to list instances: %w", err)
}

if len(instances) == 0 {
output.PrintInfo("No compute instances found")
return nil
}

dataSlice := make([]interface{}, len(instances))
for i, inst := range instances {
dataSlice[i] = inst
}

return output.Print(
dataSlice,
output.Format(outputFormat),
[]string{"UUID", "NAME", "STATE", "REGION", "IP ADDRESS"},
func(item interface{}) []string {
inst := item.(models.Instance)
return []string{inst.UUID, inst.Name, inst.State, inst.ZoneName, inst.IPAddress}
},
)
},
}

// computeGetCmd gets a specific instance
var computeGetCmd = &cobra.Command{
Use:   "get <uuid>",
Short: "Get compute instance details",
Long:  `Get detailed information about a specific compute instance.`,
Args:  cobra.ExactArgs(1),
RunE: func(cmd *cobra.Command, args []string) error {
cfg, err := config.LoadConfig()
if err != nil {
return err
}

c := client.NewClient(cfg.AccessKey, cfg.SecretKey)
uuid := args[0]

region := regionFlag
if region == "" {
region = cfg.DefaultRegion
}

instance, err := c.GetInstance(uuid, region)
if err != nil {
return fmt.Errorf("failed to get instance: %w", err)
}

return output.Print(
instance,
output.Format(outputFormat),
[]string{"UUID", "NAME", "STATE", "VCPU", "MEMORY (MB)", "DISK (GB)", "NETWORK", "PRIVATE IP", "STATUS"},
func(item interface{}) []string {
inst := item.(*models.Instance)

vcpu := inst.CPUCore
if vcpu == "" || vcpu == "0" {
vcpu = "-"
}

disk := inst.VolumeSize
if disk == "" || disk == "null" {
disk = "-"
} else {
// Converter bytes para GB (aproximado)
disk = disk + " bytes"
}

network := inst.NetworkName
if network == "" {
network = "-"
}

privateIP := inst.PrivateIP
if privateIP == "" {
privateIP = "-"
}

status := inst.Status
if status == "" {
status = inst.State
}

return []string{inst.UUID, inst.Name, inst.State, vcpu, inst.MemoryMB, disk, network, privateIP, status}
},
)
},
}

// computeCreateCmd creates a new instance
var computeCreateCmd = &cobra.Command{
Use:   "create",
Short: "Create a compute instance",
Long:  `Create a new compute instance with specified configuration.`,
RunE: func(cmd *cobra.Command, args []string) error {
cfg, err := config.LoadConfig()
if err != nil {
return err
}

name, _ := cmd.Flags().GetString("name")
region, _ := cmd.Flags().GetString("region")
templateUUID, _ := cmd.Flags().GetString("image")
offeringUUID, _ := cmd.Flags().GetString("size")
networkUUID, _ := cmd.Flags().GetString("network")
sshKey, _ := cmd.Flags().GetString("ssh-key")

if region == "" {
region = cfg.DefaultRegion
}

if name == "" || templateUUID == "" || offeringUUID == "" || networkUUID == "" {
return fmt.Errorf("required flags: --name, --image, --size, --network")
}

c := client.NewClient(cfg.AccessKey, cfg.SecretKey)

req := models.CreateInstanceRequest{
Name:                name,
Region:              region,
TemplateUUID:        templateUUID,
ComputeOfferingUUID: offeringUUID,
NetworkUUID:         networkUUID,
SSHKeyName:          sshKey,
}

output.PrintInfo(fmt.Sprintf("Creating compute instance '%s' in region '%s'...", name, region))

instance, err := c.CreateInstance(req)
if err != nil {
return fmt.Errorf("failed to create instance: %w", err)
}

output.PrintSuccess(fmt.Sprintf("Instance created: %s (UUID: %s)", instance.Name, instance.UUID))
return nil
},
}

// computeStartCmd starts an instance
var computeStartCmd = &cobra.Command{
Use:   "start <uuid>",
Short: "Start a compute instance",
Long:  `Start a stopped compute instance.`,
Args:  cobra.ExactArgs(1),
RunE: func(cmd *cobra.Command, args []string) error {
cfg, err := config.LoadConfig()
if err != nil {
return err
}

c := client.NewClient(cfg.AccessKey, cfg.SecretKey)
uuid := args[0]

output.PrintInfo(fmt.Sprintf("Starting instance %s...", uuid))

if err := c.StartInstance(uuid); err != nil {
return fmt.Errorf("failed to start instance: %w", err)
}

output.PrintSuccess(fmt.Sprintf("Instance %s started successfully", uuid))
return nil
},
}

// computeStopCmd stops an instance
var computeStopCmd = &cobra.Command{
Use:   "stop <uuid>",
Short: "Stop a compute instance",
Long:  `Stop a running compute instance.`,
Args:  cobra.ExactArgs(1),
RunE: func(cmd *cobra.Command, args []string) error {
cfg, err := config.LoadConfig()
if err != nil {
return err
}

c := client.NewClient(cfg.AccessKey, cfg.SecretKey)
uuid := args[0]

output.PrintInfo(fmt.Sprintf("Stopping instance %s...", uuid))

if err := c.StopInstance(uuid); err != nil {
return fmt.Errorf("failed to stop instance: %w", err)
}

output.PrintSuccess(fmt.Sprintf("Instance %s stopped successfully", uuid))
return nil
},
}

// computeDeleteCmd deletes an instance
var computeDeleteCmd = &cobra.Command{
Use:   "delete <uuid>",
Short: "Delete a compute instance",
Long:  `Delete a compute instance permanently. This action cannot be undone.`,
Args:  cobra.ExactArgs(1),
RunE: func(cmd *cobra.Command, args []string) error {
cfg, err := config.LoadConfig()
if err != nil {
return err
}

c := client.NewClient(cfg.AccessKey, cfg.SecretKey)
uuid := args[0]

output.PrintInfo(fmt.Sprintf("Deleting instance %s...", uuid))

if err := c.DeleteInstance(uuid); err != nil {
return fmt.Errorf("failed to delete instance: %w", err)
}

output.PrintSuccess(fmt.Sprintf("Instance %s deleted successfully", uuid))
return nil
},
}

// computeImagesCmd lists available images/templates
var computeImagesCmd = &cobra.Command{
Use:     "images",
Aliases: []string{"templates"},
Short:   "List available images",
Long:    `List all available OS images/templates that can be used to create instances.`,
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

templates, err := c.ListTemplates(region)
if err != nil {
return fmt.Errorf("failed to list images: %w", err)
}

if len(templates) == 0 {
output.PrintInfo("No images found")
return nil
}

dataSlice := make([]interface{}, len(templates))
for i, tpl := range templates {
dataSlice[i] = tpl
}

return output.Print(
dataSlice,
output.Format(outputFormat),
[]string{"UUID", "NAME", "OS TYPE", "REGION", "READY"},
func(item interface{}) []string {
tpl := item.(models.Template)
ready := "no"
if tpl.IsReady {
ready = "yes"
}
return []string{tpl.UUID, tpl.Name, tpl.OsTypeName, tpl.ZoneName, ready}
},
)
},
}

// computeSizesCmd lists available compute offerings
var computeSizesCmd = &cobra.Command{
Use:     "sizes",
Aliases: []string{"offerings", "flavors"},
Short:   "List available compute sizes",
Long:    `List all available compute offerings (VM sizes/flavors) with their specifications.`,
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

offerings, err := c.ListComputeOfferings(region)
if err != nil {
return fmt.Errorf("failed to list compute sizes: %w", err)
}

if len(offerings) == 0 {
output.PrintInfo("No compute sizes found")
return nil
}

dataSlice := make([]interface{}, len(offerings))
for i, off := range offerings {
dataSlice[i] = off
}

return output.Print(
dataSlice,
output.Format(outputFormat),
[]string{"UUID", "NAME", "CPU", "MEMORY (MB)", "ACTIVE"},
func(item interface{}) []string {
off := item.(models.ComputeOffering)
active := "no"
if off.IsActive {
active = "yes"
}
return []string{
off.UUID,
off.Name,
off.NumberOfCores,
off.Memory,
active,
}
},
)
},
}

func init() {
rootCmd.AddCommand(computeCmd)
computeCmd.AddCommand(computeListCmd)
computeCmd.AddCommand(computeGetCmd)
computeCmd.AddCommand(computeCreateCmd)
computeCmd.AddCommand(computeStartCmd)
computeCmd.AddCommand(computeStopCmd)
computeCmd.AddCommand(computeDeleteCmd)
computeCmd.AddCommand(computeImagesCmd)
computeCmd.AddCommand(computeSizesCmd)

computeCreateCmd.Flags().String("name", "", "Instance name (required)")
computeCreateCmd.Flags().String("region", "", "Region (uses default if not specified)")
computeCreateCmd.Flags().String("image", "", "Image/template UUID (required)")
computeCreateCmd.Flags().String("size", "", "Compute offering UUID (required)")
computeCreateCmd.Flags().String("network", "", "Network UUID (required)")
computeCreateCmd.Flags().String("ssh-key", "", "SSH key name")
}
