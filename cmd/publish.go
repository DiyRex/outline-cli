package cmd

import (
	"fmt"
	"os"

	"outline-cli/internal/cli"
	"outline-cli/internal/models"
	"outline-cli/internal/sync"

	"github.com/spf13/cobra"
)

var publishCmd = &cobra.Command{
	Use:   "publish <file.md>",
	Short: "Publish a markdown file to Outline (upsert)",
	Long:  "Upload a local markdown file as a document. If a document with the same title exists in the collection, it will be updated.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]
		collectionID, _ := cmd.Flags().GetString("collection")
		parentID, _ := cmd.Flags().GetString("parent")
		titleOverride, _ := cmd.Flags().GetString("title")

		if collectionID == "" {
			return fmt.Errorf("--collection is required")
		}

		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("read file: %w", err)
		}

		title := titleOverride
		if title == "" {
			title = sync.ExtractTitleFromContent(string(content), filePath)
		}

		client, err := getClient()
		if err != nil {
			return err
		}

		// Search for existing doc with same title (upsert)
		results, _, _ := client.Documents.Search(getContext(), models.SearchParams{
			Query:        title,
			CollectionID: collectionID,
		})
		for _, r := range results {
			if r.Document.Title == title && r.Document.CollectionID == collectionID {
				_, err := client.Documents.Update(getContext(), models.DocumentUpdateParams{
					ID:   r.Document.ID,
					Text: string(content),
				})
				if err != nil {
					return fmt.Errorf("update: %w", err)
				}
				cli.Output.Success("Updated: %s (%s)", title, r.Document.ID)
				return nil
			}
		}

		doc, err := client.Documents.Create(getContext(), models.DocumentCreateParams{
			Title:            title,
			Text:             string(content),
			CollectionID:     collectionID,
			ParentDocumentID: parentID,
			Publish:          true,
		})
		if err != nil {
			return err
		}
		cli.Output.Success("Published: %s (%s)", doc.Title, doc.ID)
		return nil
	},
}

func init() {
	publishCmd.Flags().String("collection", "", "Target collection ID (required)")
	publishCmd.Flags().String("parent", "", "Parent document ID")
	publishCmd.Flags().String("title", "", "Override document title")
}
