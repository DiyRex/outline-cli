package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"outline-cli/internal/cli"
	"outline-cli/internal/models"
	syncp "outline-cli/internal/sync"

	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup all collections to local folder",
	Long:  "Download all collections as markdown files organized by collection name.",
	RunE: func(cmd *cobra.Command, args []string) error {
		output, _ := cmd.Flags().GetString("output")
		if output == "" {
			output = fmt.Sprintf("outline-backup-%s", time.Now().Format("2006-01-02"))
		}

		client, err := getClient()
		if err != nil {
			return err
		}

		cli.Output.Info("Starting backup to %s/", output)

		collections, _, err := client.Collections.List(getContext(), models.CollectionListParams{
			PaginationParams: models.PaginationParams{Limit: 100},
		})
		if err != nil {
			return fmt.Errorf("list collections: %w", err)
		}

		if err := os.MkdirAll(output, 0755); err != nil {
			return fmt.Errorf("create output dir: %w", err)
		}

		totalDocs := 0
		for i, coll := range collections {
			cli.Output.Progress(i+1, len(collections), coll.Name)

			collDir := filepath.Join(output, syncp.SanitizeFilename(coll.Name))
			result, err := syncp.Pull(getContext(), client, coll.Name, collDir, syncp.PullOptions{})
			if err != nil {
				cli.Output.Warn("Failed to backup collection '%s': %s", coll.Name, err)
				continue
			}
			totalDocs += result.Downloaded
		}

		cli.Output.Success("Backup complete: %d collections, %d documents to %s/", len(collections), totalDocs, output)
		return nil
	},
}

func init() {
	backupCmd.Flags().StringP("output", "o", "", "Output directory (default: outline-backup-YYYY-MM-DD)")
}
