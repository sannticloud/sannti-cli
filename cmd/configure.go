package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/sannticloud/sannti-cli/internal/config"
	"github.com/sannticloud/sannti-cli/internal/output"
	"golang.org/x/term"
)

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure Sannti CLI credentials",
	Long: `Configure Sannti CLI with your API credentials and default settings.

This command will prompt you for:
- Sannti Access Key
- Sannti Secret Key
- Default Region

The configuration will be saved to ~/.sannti/config

You can also use environment variables:
- SANNTI_ACCESS_KEY
- SANNTI_SECRET_KEY
- SANNTI_REGION`,
	RunE: func(cmd *cobra.Command, args []string) error {
		reader := bufio.NewReader(os.Stdin)

		// Access Key
		fmt.Print("Sannti Access Key: ")
		accessKey, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read access key: %w", err)
		}
		accessKey = strings.TrimSpace(accessKey)

		// Secret Key (hidden input)
		fmt.Print("Sannti Secret Key: ")
		secretKeyBytes, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return fmt.Errorf("failed to read secret key: %w", err)
		}
		fmt.Println() // New line after hidden input
		secretKey := strings.TrimSpace(string(secretKeyBytes))

		// Default Region
		fmt.Print("Default Region [br-southeast-1]: ")
		defaultRegion, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read default region: %w", err)
		}
		defaultRegion = strings.TrimSpace(defaultRegion)
		if defaultRegion == "" {
			defaultRegion = "br-southeast-1"
		}

		// Validate credentials are not empty
		if accessKey == "" || secretKey == "" {
			return fmt.Errorf("access key and secret key cannot be empty")
		}

		// Save configuration
		if err := config.SaveConfig(accessKey, secretKey, defaultRegion); err != nil {
			return fmt.Errorf("failed to save configuration: %w", err)
		}

		configPath, _ := config.GetConfigPath()
		output.PrintSuccess(fmt.Sprintf("Configuration saved to %s/config.yaml", configPath))
		output.PrintInfo("You can now use Sannti CLI commands!")
		output.PrintInfo("Try: sannti region list")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(configureCmd)
}
