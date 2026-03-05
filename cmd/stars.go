package cmd

import (
	"fmt"
	"outline-cli/internal/models"

	"github.com/spf13/cobra"
)

var starsCmd = &cobra.Command{
	Use:   "stars",
	Short: "Manage starred items",
}

var starsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List starred items",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		stars, _, err := client.Stars.List(getContext(), models.StarListParams{
			PaginationParams: models.PaginationParams{Limit: 100},
		})
		if err != nil {
			return err
		}
		var rows [][]string
		for _, s := range stars {
			target := s.DocumentID
			if target == "" {
				target = s.CollectionID
			}
			rows = append(rows, []string{s.ID, target, s.CreatedAt.Format("2006-01-02 15:04")})
		}
		printOutput(stars, []string{"ID", "TARGET", "CREATED"}, rows)
		return nil
	},
}

var starsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Star a document or collection",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		docID, _ := cmd.Flags().GetString("document")
		collID, _ := cmd.Flags().GetString("collection")
		star, err := client.Stars.Create(getContext(), models.StarCreateParams{
			DocumentID:   docID,
			CollectionID: collID,
		})
		if err != nil {
			return err
		}
		fmt.Printf("Starred: %s\n", star.ID)
		return nil
	},
}

var starsDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Unstar an item",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		if err := client.Stars.Delete(getContext(), args[0]); err != nil {
			return err
		}
		fmt.Println("Unstarred.")
		return nil
	},
}

func init() {
	starsCmd.AddCommand(starsListCmd)
	starsCmd.AddCommand(starsCreateCmd)
	starsCmd.AddCommand(starsDeleteCmd)

	starsCreateCmd.Flags().String("document", "", "Document ID")
	starsCreateCmd.Flags().String("collection", "", "Collection ID")
}
