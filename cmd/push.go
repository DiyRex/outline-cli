package cmd

import (
	"fmt"

	"outline-cli/internal/sync"

	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:   "push <collection-name> <folder-path>",
	Short: "Push a local folder as an Outline collection",
	Long:  "Upload a folder of markdown files to Outline, preserving directory hierarchy as nested documents.",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}

		collectionName := args[0]
		folderPath := args[1]
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		del, _ := cmd.Flags().GetBool("delete")

		opts := sync.PushOptions{
			DryRun: dryRun,
			Delete: del,
		}

		result, err := sync.Push(getContext(), client, collectionName, folderPath, opts)
		if err != nil {
			return err
		}

		fmt.Printf("\nPush complete: %d created, %d updated, %d skipped\n", result.Created, result.Updated, result.Skipped)
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
	pushCmd.Flags().Bool("dry-run", false, "Preview changes without uploading")
	pushCmd.Flags().Bool("delete", false, "Delete remote documents not present locally")
}
