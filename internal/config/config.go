package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sardonyx001/unlinked/pkg/types"
	"github.com/spf13/viper"
)

// Manager handles application configuration
type Manager struct {
	v      *viper.Viper
	config *types.Config
}

// New creates a new configuration manager
func New() *Manager {
	return &Manager{
		v:      viper.New(),
		config: types.DefaultConfig(),
	}
}

// Load loads configuration from file and environment
func (m *Manager) Load(configFile string) error {
	// Set config file if provided
	if configFile != "" {
		m.v.SetConfigFile(configFile)
	} else {
		// Search for config in common locations
		home, err := os.UserHomeDir()
		if err == nil {
			m.v.AddConfigPath(filepath.Join(home, ".config", "unlinked"))
		}
		m.v.AddConfigPath(".")
		m.v.SetConfigName("config")
		m.v.SetConfigType("yaml")
	}

	// Environment variable support
	m.v.SetEnvPrefix("UNLINKED")
	m.v.AutomaticEnv()

	// Set defaults
	m.setDefaults()

	// Read config file (optional)
	if err := m.v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found is OK, we'll use defaults
	}

	// Unmarshal into config struct
	if err := m.v.Unmarshal(m.config); err != nil {
		return fmt.Errorf("error unmarshaling config: %w", err)
	}

	return nil
}

// setDefaults sets default values for all configuration options
func (m *Manager) setDefaults() {
	defaults := types.DefaultConfig()

	m.v.SetDefault("mode", defaults.Mode)
	m.v.SetDefault("output_format", defaults.OutputFormat)
	m.v.SetDefault("concurrency", defaults.Concurrency)
	m.v.SetDefault("timeout", defaults.Timeout)
	m.v.SetDefault("max_depth", defaults.MaxDepth)
	m.v.SetDefault("follow_redirects", defaults.FollowRedirects)
	m.v.SetDefault("check_external_only", defaults.CheckExternalOnly)
	m.v.SetDefault("user_agent", defaults.UserAgent)
	m.v.SetDefault("respect_robots_txt", defaults.RespectRobotsTxt)
	m.v.SetDefault("verbose", defaults.Verbose)
	m.v.SetDefault("show_progress", defaults.ShowProgress)
}

// Get returns the current configuration
func (m *Manager) Get() *types.Config {
	return m.config
}

// Set updates a configuration value
func (m *Manager) Set(key string, value interface{}) {
	m.v.Set(key, value)
	// Re-unmarshal to update config struct
	m.v.Unmarshal(m.config)
}

// GetViper returns the underlying viper instance
func (m *Manager) GetViper() *viper.Viper {
	return m.v
}
