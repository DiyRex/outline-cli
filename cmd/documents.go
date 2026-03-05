package cmd

import (
	"fmt"
	"io"
	"os"
	"outline-cli/internal/cli"
	"outline-cli/internal/models"

	"github.com/spf13/cobra"
)

var documentsCmd = &cobra.Command{
	Use:     "documents",
	Short:   "Manage documents",
	Aliases: []string{"doc", "docs"},
}

var documentsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List documents",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		collectionID, _ := cmd.Flags().GetString("collection")
		docs, _, err := client.Documents.List(getContext(), models.DocumentListParams{
			CollectionID:     collectionID,
			PaginationParams: models.PaginationParams{Limit: 100},
		})
		if err != nil {
			return err
		}
		var rows [][]string
		for _, d := range docs {
			rows = append(rows, []string{d.ID, d.Title, d.UpdatedAt.Format("2006-01-02 15:04")})
		}
		printOutput(docs, []string{"ID", "TITLE", "UPDATED"}, rows)
		return nil
	},
}

var documentsInfoCmd = &cobra.Command{
	Use:   "info <id>",
	Short: "Get document details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		doc, err := client.Documents.Info(getContext(), args[0])
		if err != nil {
			return err
		}
		if outputFormat == "json" {
			printJSON(doc)
		} else {
			fmt.Printf("ID:         %s\n", doc.ID)
			fmt.Printf("Title:      %s\n", doc.Title)
			fmt.Printf("Collection: %s\n", doc.CollectionID)
			fmt.Printf("Updated:    %s\n", doc.UpdatedAt.Format("2006-01-02 15:04"))
			fmt.Printf("Revision:   %d\n", doc.Revision)
			fmt.Println("\n--- Content ---")
			fmt.Println(doc.Text)
		}
		return nil
	},
}

var documentsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a document",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		title, _ := cmd.Flags().GetString("title")
		text, _ := cmd.Flags().GetString("text")
		collectionID, _ := cmd.Flags().GetString("collection")
		parentID, _ := cmd.Flags().GetString("parent")
		file, _ := cmd.Flags().GetString("file")
		useStdin, _ := cmd.Flags().GetBool("stdin")
		publish, _ := cmd.Flags().GetBool("publish")
		templateID, _ := cmd.Flags().GetString("template")

		if useStdin {
			data, err := io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("read stdin: %w", err)
			}
			text = string(data)
		} else if file != "" {
			content, err := os.ReadFile(file)
			if err != nil {
				return fmt.Errorf("read file: %w", err)
			}
			text = string(content)
		}

		doc, err := client.Documents.Create(getContext(), models.DocumentCreateParams{
			Title:            title,
			Text:             text,
			CollectionID:     collectionID,
			ParentDocumentID: parentID,
			Publish:          publish,
			TemplateID:       templateID,
		})
		if err != nil {
			return err
		}
		fmt.Printf("Created document: %s (%s)\n", doc.Title, doc.ID)
		return nil
	},
}

var documentsUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a document",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		title, _ := cmd.Flags().GetString("title")
		text, _ := cmd.Flags().GetString("text")
		file, _ := cmd.Flags().GetString("file")
		useStdin, _ := cmd.Flags().GetBool("stdin")

		if useStdin {
			data, err := io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("read stdin: %w", err)
			}
			text = string(data)
		} else if file != "" {
			content, err := os.ReadFile(file)
			if err != nil {
				return fmt.Errorf("read file: %w", err)
			}
			text = string(content)
		}

		doc, err := client.Documents.Update(getContext(), models.DocumentUpdateParams{
			ID:    args[0],
			Title: title,
			Text:  text,
		})
		if err != nil {
			return err
		}
		fmt.Printf("Updated document: %s (%s)\n", doc.Title, doc.ID)
		return nil
	},
}

var documentsDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a document",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		permanent, _ := cmd.Flags().GetBool("permanent")
		if permanent && !cli.Output.Confirm("Permanently delete this document? This cannot be undone") {
			fmt.Println("Cancelled.")
			return nil
		}
		client, err := getClient()
		if err != nil {
			return err
		}
		if err := client.Documents.Delete(getContext(), args[0], permanent); err != nil {
			return err
		}
		fmt.Println("Document deleted.")
		return nil
	},
}

var documentsArchiveCmd = &cobra.Command{
	Use:   "archive <id>",
	Short: "Archive a document",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		doc, err := client.Documents.Archive(getContext(), args[0])
		if err != nil {
			return err
		}
		fmt.Printf("Archived: %s\n", doc.Title)
		return nil
	},
}

var documentsRestoreCmd = &cobra.Command{
	Use:   "restore <id>",
	Short: "Restore a document",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		doc, err := client.Documents.Restore(getContext(), args[0])
		if err != nil {
			return err
		}
		fmt.Printf("Restored: %s\n", doc.Title)
		return nil
	},
}

var documentsMoveCmd = &cobra.Command{
	Use:   "move <id>",
	Short: "Move a document",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		collectionID, _ := cmd.Flags().GetString("collection")
		parentID, _ := cmd.Flags().GetString("parent")
		doc, err := client.Documents.Move(getContext(), models.DocumentMoveParams{
			ID:               args[0],
			CollectionID:     collectionID,
			ParentDocumentID: parentID,
		})
		if err != nil {
			return err
		}
		fmt.Printf("Moved: %s\n", doc.Title)
		return nil
	},
}

var documentsExportCmd = &cobra.Command{
	Use:   "export <id>",
	Short: "Export document content",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		output, _ := cmd.Flags().GetString("output")
		doc, err := client.Documents.Info(getContext(), args[0])
		if err != nil {
			return err
		}
		if output != "" {
			if err := os.WriteFile(output, []byte(doc.Text), 0644); err != nil {
				return err
			}
			fmt.Printf("Exported to: %s\n", output)
		} else {
			fmt.Println(doc.Text)
		}
		return nil
	},
}

var documentsDuplicateCmd = &cobra.Command{
	Use:   "duplicate <id>",
	Short: "Duplicate a document",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		recursive, _ := cmd.Flags().GetBool("recursive")
		doc, err := client.Documents.Duplicate(getContext(), args[0], recursive)
		if err != nil {
			return err
		}
		fmt.Printf("Duplicated: %s (%s)\n", doc.Title, doc.ID)
		return nil
	},
}

var documentsSearchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search documents",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		collectionID, _ := cmd.Flags().GetString("collection")
		results, _, err := client.Documents.Search(getContext(), models.SearchParams{
			Query:        args[0],
			CollectionID: collectionID,
		})
		if err != nil {
			return err
		}
		var rows [][]string
		for _, r := range results {
			rows = append(rows, []string{r.Document.ID, r.Document.Title, fmt.Sprintf("%.2f", r.Ranking)})
		}
		printOutput(results, []string{"ID", "TITLE", "SCORE"}, rows)
		return nil
	},
}

var documentsDraftsCmd = &cobra.Command{
	Use:   "drafts",
	Short: "List unpublished drafts",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		collectionID, _ := cmd.Flags().GetString("collection")
		docs, _, err := client.Documents.Drafts(getContext(), models.DocumentListParams{
			CollectionID:     collectionID,
			PaginationParams: models.PaginationParams{Limit: 100},
		})
		if err != nil {
			return err
		}
		var rows [][]string
		for _, d := range docs {
			rows = append(rows, []string{d.ID, d.Title, d.UpdatedAt.Format("2006-01-02 15:04")})
		}
		printOutput(docs, []string{"ID", "TITLE", "UPDATED"}, rows)
		return nil
	},
}

var documentsViewedCmd = &cobra.Command{
	Use:   "viewed",
	Short: "List recently viewed documents",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		docs, _, err := client.Documents.Viewed(getContext(), models.PaginationParams{Limit: 50})
		if err != nil {
			return err
		}
		var rows [][]string
		for _, d := range docs {
			rows = append(rows, []string{d.ID, d.Title, d.UpdatedAt.Format("2006-01-02 15:04")})
		}
		printOutput(docs, []string{"ID", "TITLE", "UPDATED"}, rows)
		return nil
	},
}

var documentsUnpublishCmd = &cobra.Command{
	Use:   "unpublish <id>",
	Short: "Unpublish a document (return to draft)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		doc, err := client.Documents.Unpublish(getContext(), args[0])
		if err != nil {
			return err
		}
		fmt.Printf("Unpublished: %s\n", doc.Title)
		return nil
	},
}

func init() {
	documentsCmd.AddCommand(documentsListCmd)
	documentsCmd.AddCommand(documentsInfoCmd)
	documentsCmd.AddCommand(documentsCreateCmd)
	documentsCmd.AddCommand(documentsUpdateCmd)
	documentsCmd.AddCommand(documentsDeleteCmd)
	documentsCmd.AddCommand(documentsArchiveCmd)
	documentsCmd.AddCommand(documentsRestoreCmd)
	documentsCmd.AddCommand(documentsMoveCmd)
	documentsCmd.AddCommand(documentsExportCmd)
	documentsCmd.AddCommand(documentsDuplicateCmd)
	documentsCmd.AddCommand(documentsSearchCmd)
	documentsCmd.AddCommand(documentsDraftsCmd)
	documentsCmd.AddCommand(documentsViewedCmd)
	documentsCmd.AddCommand(documentsUnpublishCmd)

	documentsListCmd.Flags().String("collection", "", "Filter by collection ID")
	documentsCreateCmd.Flags().String("title", "", "Document title")
	documentsCreateCmd.Flags().String("text", "", "Document content (markdown)")
	documentsCreateCmd.Flags().String("collection", "", "Collection ID")
	documentsCreateCmd.Flags().String("parent", "", "Parent document ID")
	documentsCreateCmd.Flags().String("file", "", "Read content from file")
	documentsCreateCmd.Flags().String("template", "", "Template ID to create from")
	documentsCreateCmd.Flags().Bool("publish", true, "Publish immediately")
	documentsCreateCmd.Flags().Bool("stdin", false, "Read content from stdin")
	documentsUpdateCmd.Flags().String("title", "", "New title")
	documentsUpdateCmd.Flags().String("text", "", "New content")
	documentsUpdateCmd.Flags().String("file", "", "Read content from file")
	documentsUpdateCmd.Flags().Bool("stdin", false, "Read content from stdin")
	documentsDeleteCmd.Flags().Bool("permanent", false, "Permanently delete")
	documentsMoveCmd.Flags().String("collection", "", "Target collection ID")
	documentsMoveCmd.Flags().String("parent", "", "Target parent document ID")
	documentsExportCmd.Flags().String("output", "", "Output file path")
	documentsDuplicateCmd.Flags().Bool("recursive", false, "Include child documents")
	documentsSearchCmd.Flags().String("collection", "", "Filter by collection ID")
	documentsDraftsCmd.Flags().String("collection", "", "Filter by collection ID")
}
