package cmd

import (
	"fmt"
	"outline-cli/internal/models"

	"github.com/spf13/cobra"
)

var commentsCmd = &cobra.Command{
	Use:   "comments",
	Short: "Manage comments",
}

var commentsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List comments",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		docID, _ := cmd.Flags().GetString("document")
		comments, _, err := client.Comments.List(getContext(), models.CommentListParams{
			DocumentID:       docID,
			PaginationParams: models.PaginationParams{Limit: 100},
		})
		if err != nil {
			return err
		}
		var rows [][]string
		for _, c := range comments {
			createdBy := ""
			if c.CreatedBy != nil {
				createdBy = c.CreatedBy.Name
			}
			rows = append(rows, []string{c.ID, c.DocumentID, createdBy, c.CreatedAt.Format("2006-01-02 15:04")})
		}
		printOutput(comments, []string{"ID", "DOCUMENT", "AUTHOR", "CREATED"}, rows)
		return nil
	},
}

var commentsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a comment",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		docID, _ := cmd.Flags().GetString("document")
		text, _ := cmd.Flags().GetString("text")
		parentID, _ := cmd.Flags().GetString("parent")
		comment, err := client.Comments.Create(getContext(), models.CommentCreateParams{
			DocumentID:      docID,
			Text:            text,
			ParentCommentID: parentID,
		})
		if err != nil {
			return err
		}
		fmt.Printf("Created comment: %s\n", comment.ID)
		return nil
	},
}

var commentsDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a comment",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		if err := client.Comments.Delete(getContext(), args[0]); err != nil {
			return err
		}
		fmt.Println("Comment deleted.")
		return nil
	},
}

var commentsResolveCmd = &cobra.Command{
	Use:   "resolve <id>",
	Short: "Resolve a comment thread",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		_, err = client.Comments.Resolve(getContext(), args[0])
		if err != nil {
			return err
		}
		fmt.Println("Comment resolved.")
		return nil
	},
}

func init() {
	commentsCmd.AddCommand(commentsListCmd)
	commentsCmd.AddCommand(commentsCreateCmd)
	commentsCmd.AddCommand(commentsDeleteCmd)
	commentsCmd.AddCommand(commentsResolveCmd)

	commentsListCmd.Flags().String("document", "", "Filter by document ID")
	commentsCreateCmd.Flags().String("document", "", "Document ID")
	commentsCreateCmd.Flags().String("text", "", "Comment text")
	commentsCreateCmd.Flags().String("parent", "", "Parent comment ID")
}
