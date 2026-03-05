package cmd

import (
	"fmt"

	"outline-cli/internal/config"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure Outline CLI",
	Long:  "Set the Outline instance URL and API key.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			cfg = &config.Config{}
		}

		url, _ := cmd.Flags().GetString("url")
		apiKey, _ := cmd.Flags().GetString("api-key")

		if url == "" && apiKey == "" {
			fmt.Printf("Config file: %s\n", config.ConfigPath())
			fmt.Printf("URL:         %s\n", cfg.URL)
			if cfg.APIKey != "" {
				fmt.Printf("API Key:     %s...%s\n", cfg.APIKey[:10], cfg.APIKey[len(cfg.APIKey)-4:])
			} else {
				fmt.Printf("API Key:     (not set)\n")
			}
			return nil
		}

		if url != "" {
			cfg.URL = url
		}
		if apiKey != "" {
			cfg.APIKey = apiKey
		}

		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("save config: %w", err)
		}

		fmt.Println("Configuration saved.")
		return nil
	},
}

func init() {
	configCmd.Flags().String("url", "", "Outline instance URL (e.g., https://outline.example.com)")
	configCmd.Flags().String("api-key", "", "Outline API key (starts with ol_api_)")
}
