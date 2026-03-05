package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"outline-cli/internal/api"
	"outline-cli/internal/cli"
	"outline-cli/internal/config"

	"github.com/spf13/cobra"
)

var (
	outputFormat   string
	quietMode      bool
	verboseMode    bool
	noColor        bool
	nonInteractive bool
	yesMode        bool
	timeoutSecs    int

	rootCmd = &cobra.Command{
		Use:   "outline",
		Short: "CLI tool for Outline wiki",
		Long:  "A comprehensive CLI and TUI for managing Outline wiki instances.",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			cli.Output.Quiet = quietMode
			cli.Output.Verbose = verboseMode
			cli.Output.NoColor = noColor || os.Getenv("NO_COLOR") != ""
			cli.Output.YesMode = yesMode || nonInteractive
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if nonInteractive {
				return fmt.Errorf("no subcommand specified. Use --help to see available commands")
			}
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			if err := cfg.Validate(); err != nil {
				return err
			}
			client := getClientFromConfig(cfg)
			return runTUI(client)
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		code := cli.ExitCodeFromAPIError(err)
		os.Exit(code)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "format", "f", "table", "Output format: table or json")
	rootCmd.PersistentFlags().BoolVarP(&quietMode, "quiet", "q", false, "Suppress non-error output")
	rootCmd.PersistentFlags().BoolVarP(&verboseMode, "verbose", "v", false, "Enable debug output")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "Disable colored output")
	rootCmd.PersistentFlags().BoolVar(&nonInteractive, "non-interactive", false, "Disable interactive prompts and TUI")
	rootCmd.PersistentFlags().BoolVarP(&yesMode, "yes", "y", false, "Auto-confirm destructive operations")
	rootCmd.PersistentFlags().IntVar(&timeoutSecs, "timeout", 30, "HTTP timeout in seconds")

	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(pushCmd)
	rootCmd.AddCommand(pullCmd)
	rootCmd.AddCommand(documentsCmd)
	rootCmd.AddCommand(collectionsCmd)
	rootCmd.AddCommand(usersCmd)
	rootCmd.AddCommand(groupsCmd)
	rootCmd.AddCommand(commentsCmd)
	rootCmd.AddCommand(sharesCmd)
	rootCmd.AddCommand(starsCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(eventsCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(changelogCmd)
	rootCmd.AddCommand(publishCmd)
	rootCmd.AddCommand(revisionsCmd)
	rootCmd.AddCommand(attachmentsCmd)
	rootCmd.AddCommand(diffCmd)
	rootCmd.AddCommand(backupCmd)
}

func getClientFromConfig(cfg *config.Config) *api.Client {
	client := api.NewClient(cfg.URL, cfg.APIKey)
	if timeoutSecs > 0 {
		client.HTTPClient.Timeout = time.Duration(timeoutSecs) * time.Second
	}
	return client
}

func getClient() (*api.Client, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return getClientFromConfig(cfg), nil
}

func getContext() context.Context {
	return context.Background()
}

func printJSON(v interface{}) {
	data, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(data))
}

func newTabWriter() *tabwriter.Writer {
	return tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
}

func printOutput(v interface{}, headers []string, rows [][]string) {
	if outputFormat == "json" {
		printJSON(v)
		return
	}
	w := newTabWriter()
	header := ""
	for i, h := range headers {
		if i > 0 {
			header += "\t"
		}
		header += h
	}
	fmt.Fprintln(w, header)
	for _, row := range rows {
		line := ""
		for i, col := range row {
			if i > 0 {
				line += "\t"
			}
			line += col
		}
		fmt.Fprintln(w, line)
	}
	w.Flush()
}
