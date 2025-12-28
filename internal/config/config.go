package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	ConfigFileName = "config"
	ConfigFileType = "yaml"
	ConfigDirName  = ".sannti"
)

// Config holds the CLI configuration
type Config struct {
	AccessKey     string `mapstructure:"access_key"`
	SecretKey     string `mapstructure:"secret_key"`
	DefaultRegion string `mapstructure:"default_region"`
}

// GetConfigPath returns the path to the config file
func GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(home, ConfigDirName)
	return configDir, nil
}

// InitConfig initializes the configuration
func InitConfig() error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	viper.AddConfigPath(configPath)
	viper.SetConfigName(ConfigFileName)
	viper.SetConfigType(ConfigFileType)

	// Set environment variable prefix
	viper.SetEnvPrefix("SANNTI")
	viper.AutomaticEnv()

	// Read config if exists
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("failed to read config file: %w", err)
		}
	}

	return nil
}

// LoadConfig loads the configuration with environment variable overrides
func LoadConfig() (*Config, error) {
	if err := InitConfig(); err != nil {
		return nil, err
	}

	cfg := &Config{
		AccessKey:     viper.GetString("access_key"),
		SecretKey:     viper.GetString("secret_key"),
		DefaultRegion: viper.GetString("default_region"),
	}

	// Check for missing required fields
	if cfg.AccessKey == "" || cfg.SecretKey == "" {
		return nil, fmt.Errorf("missing credentials. Please run 'sannti configure' first")
	}

	return cfg, nil
}

// SaveConfig saves the configuration to file
func SaveConfig(accessKey, secretKey, defaultRegion string) error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configPath, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Set values
	viper.Set("access_key", accessKey)
	viper.Set("secret_key", secretKey)
	viper.Set("default_region", defaultRegion)

	// Write config file
	configFile := filepath.Join(configPath, ConfigFileName+"."+ConfigFileType)
	if err := viper.WriteConfigAs(configFile); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	// Set file permissions to 0600 (read/write for owner only)
	if err := os.Chmod(configFile, 0600); err != nil {
		return fmt.Errorf("failed to set config file permissions: %w", err)
	}

	return nil
}

// GetDefaultRegion returns the default region from config or env
func GetDefaultRegion() string {
	return viper.GetString("default_region")
}
