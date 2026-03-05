package cmd

import (
	"fmt"
	"outline-cli/internal/models"

	"github.com/spf13/cobra"
)

var sharesCmd = &cobra.Command{
	Use:   "shares",
	Short: "Manage share links",
}

var sharesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List shares",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		shares, _, err := client.Shares.List(getContext(), models.ShareListParams{
			PaginationParams: models.PaginationParams{Limit: 100},
		})
		if err != nil {
			return err
		}
		var rows [][]string
		for _, s := range shares {
			rows = append(rows, []string{s.ID, s.DocumentID, s.URL, fmt.Sprintf("%v", s.Published)})
		}
		printOutput(shares, []string{"ID", "DOCUMENT", "URL", "PUBLISHED"}, rows)
		return nil
	},
}

var sharesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a share link",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		docID, _ := cmd.Flags().GetString("document")
		share, err := client.Shares.Create(getContext(), models.ShareCreateParams{
			DocumentID: docID,
		})
		if err != nil {
			return err
		}
		fmt.Printf("Share created: %s\nURL: %s\n", share.ID, share.URL)
		return nil
	},
}

var sharesRevokeCmd = &cobra.Command{
	Use:   "revoke <id>",
	Short: "Revoke a share link",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		if err := client.Shares.Revoke(getContext(), args[0]); err != nil {
			return err
		}
		fmt.Println("Share revoked.")
		return nil
	},
}

func init() {
	sharesCmd.AddCommand(sharesListCmd)
	sharesCmd.AddCommand(sharesCreateCmd)
	sharesCmd.AddCommand(sharesRevokeCmd)

	sharesCreateCmd.Flags().String("document", "", "Document ID to share")
}
