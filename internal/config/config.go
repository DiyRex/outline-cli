package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	URL    string `mapstructure:"url" yaml:"url"`
	APIKey string `mapstructure:"api_key" yaml:"api_key"`
}

func ConfigDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".outline-cli")
}

func ConfigPath() string {
	return filepath.Join(ConfigDir(), "config.yaml")
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(ConfigDir())
	viper.SetEnvPrefix("OUTLINE")
	viper.AutomaticEnv()

	_ = viper.ReadInConfig()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	return &cfg, nil
}

func Save(cfg *Config) error {
	dir := ConfigDir()
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}
	viper.Set("url", cfg.URL)
	viper.Set("api_key", cfg.APIKey)
	return viper.WriteConfigAs(ConfigPath())
}

func (c *Config) Validate() error {
	if c.URL == "" {
		return fmt.Errorf("URL not configured. Run: outline config --url=<your-outline-url>")
	}
	if c.APIKey == "" {
		return fmt.Errorf("API key not configured. Run: outline config --api-key=<your-api-key>")
	}
	return nil
}
