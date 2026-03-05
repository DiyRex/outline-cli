package cmd

import (
	"fmt"
	"outline-cli/internal/models"

	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search documents",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		collectionID, _ := cmd.Flags().GetString("collection")
		titlesOnly, _ := cmd.Flags().GetBool("titles")

		if titlesOnly {
			docs, _, err := client.Search.Titles(getContext(), models.SearchParams{
				Query:        args[0],
				CollectionID: collectionID,
			})
			if err != nil {
				return err
			}
			var rows [][]string
			for _, d := range docs {
				rows = append(rows, []string{d.ID, d.Title})
			}
			printOutput(docs, []string{"ID", "TITLE"}, rows)
		} else {
			results, _, err := client.Search.Documents(getContext(), models.SearchParams{
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
		}
		return nil
	},
}

func init() {
	searchCmd.Flags().String("collection", "", "Filter by collection ID")
	searchCmd.Flags().Bool("titles", false, "Search titles only (faster)")
}
