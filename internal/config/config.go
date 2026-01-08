package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	APIKey            string   `mapstructure:"api_key"`
	DefaultSources    []string `mapstructure:"default_sources"`
	DefaultCabins     []string `mapstructure:"default_cabins"`
	PreferredAirports []string `mapstructure:"preferred_airports"`
}

// Load reads the configuration from file and environment variables
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Config file locations
	configDir, err := os.UserConfigDir()
	if err == nil {
		viper.AddConfigPath(filepath.Join(configDir, "seats-aero"))
	}
	viper.AddConfigPath("$HOME/.config/seats-aero")
	viper.AddConfigPath(".")

	// Environment variables
	viper.SetEnvPrefix("SEATS_AERO")
	viper.BindEnv("api_key", "SEATS_AERO_API_KEY")
	viper.AutomaticEnv()

	// Defaults
	viper.SetDefault("default_sources", []string{})
	viper.SetDefault("default_cabins", []string{"J", "F"})
	viper.SetDefault("preferred_airports", []string{})

	// Read config file (ignore if not found)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &cfg, nil
}

// GetAPIKey returns the API key, checking env var first then config
func (c *Config) GetAPIKey() string {
	if envKey := os.Getenv("SEATS_AERO_API_KEY"); envKey != "" {
		return envKey
	}
	return c.APIKey
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.GetAPIKey() == "" {
		return fmt.Errorf("API key not configured. Set SEATS_AERO_API_KEY environment variable or add api_key to config file")
	}
	return nil
}

// ConfigPath returns the path where the config file should be stored
func ConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "seats-aero", "config.yaml"), nil
}

// EnsureConfigDir creates the config directory if it doesn't exist
func EnsureConfigDir() error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	dir := filepath.Join(configDir, "seats-aero")
	return os.MkdirAll(dir, 0755)
}
