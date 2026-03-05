package cmd

import (
	"fmt"
	"outline-cli/internal/models"

	"github.com/spf13/cobra"
)

var revisionsCmd = &cobra.Command{
	Use:   "revisions",
	Short: "Manage document revisions",
	Aliases: []string{"rev"},
}

var revisionsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List revisions for a document",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		docID, _ := cmd.Flags().GetString("document")
		if docID == "" {
			return fmt.Errorf("--document is required")
		}
		revisions, _, err := client.Revisions.List(getContext(), models.RevisionListParams{
			DocumentID: docID,
			PaginationParams: models.PaginationParams{Limit: 50},
		})
		if err != nil {
			return err
		}
		var rows [][]string
		for _, r := range revisions {
			author := ""
			if r.CreatedBy != nil {
				author = r.CreatedBy.Name
			}
			rows = append(rows, []string{r.ID, r.Title, fmt.Sprintf("v%d", r.Version), author, r.CreatedAt.Format("2006-01-02 15:04")})
		}
		printOutput(revisions, []string{"ID", "TITLE", "VERSION", "AUTHOR", "CREATED"}, rows)
		return nil
	},
}

var revisionsInfoCmd = &cobra.Command{
	Use:   "info <id>",
	Short: "Show revision details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		rev, err := client.Revisions.Info(getContext(), args[0])
		if err != nil {
			return err
		}
		if outputFormat == "json" {
			printJSON(rev)
		} else {
			fmt.Printf("ID:      %s\n", rev.ID)
			fmt.Printf("Title:   %s\n", rev.Title)
			fmt.Printf("Version: %d\n", rev.Version)
			fmt.Printf("Created: %s\n", rev.CreatedAt.Format("2006-01-02 15:04"))
			fmt.Println("\n--- Content ---")
			fmt.Println(rev.Text)
		}
		return nil
	},
}

var revisionsDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a revision",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		if err := client.Revisions.Delete(getContext(), args[0]); err != nil {
			return err
		}
		fmt.Println("Revision deleted.")
		return nil
	},
}

func init() {
	revisionsCmd.AddCommand(revisionsListCmd)
	revisionsCmd.AddCommand(revisionsInfoCmd)
	revisionsCmd.AddCommand(revisionsDeleteCmd)

	revisionsListCmd.Flags().String("document", "", "Document ID (required)")
}
