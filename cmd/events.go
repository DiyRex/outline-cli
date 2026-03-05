package cmd

import (
	"outline-cli/internal/models"

	"github.com/spf13/cobra"
)

var eventsCmd = &cobra.Command{
	Use:   "events",
	Short: "List activity events",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		docID, _ := cmd.Flags().GetString("document")
		collID, _ := cmd.Flags().GetString("collection")
		auditLog, _ := cmd.Flags().GetBool("audit")

		events, _, err := client.Events.List(getContext(), models.EventListParams{
			DocumentID:       docID,
			CollectionID:     collID,
			AuditLog:         auditLog,
			PaginationParams: models.PaginationParams{Limit: 50},
		})
		if err != nil {
			return err
		}
		var rows [][]string
		for _, e := range events {
			actor := ""
			if e.Actor != nil {
				actor = e.Actor.Name
			}
			rows = append(rows, []string{e.ID, e.Name, actor, e.CreatedAt.Format("2006-01-02 15:04")})
		}
		printOutput(events, []string{"ID", "EVENT", "ACTOR", "CREATED"}, rows)
		return nil
	},
}

func init() {
	eventsCmd.Flags().String("document", "", "Filter by document ID")
	eventsCmd.Flags().String("collection", "", "Filter by collection ID")
	eventsCmd.Flags().Bool("audit", false, "Show audit log events")
}
