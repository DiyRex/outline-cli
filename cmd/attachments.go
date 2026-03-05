package cmd

import (
	"fmt"
	"outline-cli/internal/models"

	"github.com/spf13/cobra"
)

var attachmentsCmd = &cobra.Command{
	Use:   "attachments",
	Short: "Manage attachments",
	Aliases: []string{"att"},
}

var attachmentsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List attachments",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		docID, _ := cmd.Flags().GetString("document")
		attachments, _, err := client.Attachments.List(getContext(), models.AttachmentListParams{
			DocumentID:       docID,
			PaginationParams: models.PaginationParams{Limit: 100},
		})
		if err != nil {
			return err
		}
		var rows [][]string
		for _, a := range attachments {
			rows = append(rows, []string{a.ID, a.Name, a.ContentType, fmt.Sprintf("%d", a.Size)})
		}
		printOutput(attachments, []string{"ID", "NAME", "TYPE", "SIZE"}, rows)
		return nil
	},
}

var attachmentsDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete an attachment",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		if err := client.Attachments.Delete(getContext(), args[0]); err != nil {
			return err
		}
		fmt.Println("Attachment deleted.")
		return nil
	},
}

func init() {
	attachmentsCmd.AddCommand(attachmentsListCmd)
	attachmentsCmd.AddCommand(attachmentsDeleteCmd)

	attachmentsListCmd.Flags().String("document", "", "Filter by document ID")
}
