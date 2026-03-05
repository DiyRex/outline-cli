package cmd

import (
	"fmt"

	"outline-cli/internal/sync"

	"github.com/spf13/cobra"
)

var pullCmd = &cobra.Command{
	Use:   "pull <collection-name> [local-path]",
	Short: "Pull an Outline collection to local folder",
	Long:  "Download a collection from Outline as markdown files, preserving document hierarchy.",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}

		collectionName := args[0]
		localPath := "."
		if len(args) > 1 {
			localPath = args[1]
		}

		dryRun, _ := cmd.Flags().GetBool("dry-run")
		opts := sync.PullOptions{DryRun: dryRun}

		result, err := sync.Pull(getContext(), client, collectionName, localPath, opts)
		if err != nil {
			return err
		}

		fmt.Printf("\nPull complete: %d documents downloaded\n", result.Downloaded)
		if len(result.Errors) > 0 {
			fmt.Printf("Errors: %d\n", len(result.Errors))
			for _, e := range result.Errors {
				fmt.Printf("  - %s\n", e)
			}
		}
		return nil
	},
}

func init() {
	pullCmd.Flags().Bool("dry-run", false, "Preview changes without downloading")
}
