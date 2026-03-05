package cmd

import (
	"fmt"
	"outline-cli/internal/cli"
	"outline-cli/internal/models"
	"strings"

	"github.com/spf13/cobra"
)

var collectionsCmd = &cobra.Command{
	Use:     "collections",
	Short:   "Manage collections",
	Aliases: []string{"coll", "col"},
}

var collectionsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List collections",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		collections, _, err := client.Collections.List(getContext(), models.CollectionListParams{
			PaginationParams: models.PaginationParams{Limit: 100},
		})
		if err != nil {
			return err
		}
		var rows [][]string
		for _, c := range collections {
			rows = append(rows, []string{c.ID, c.Name, c.UpdatedAt.Format("2006-01-02 15:04")})
		}
		printOutput(collections, []string{"ID", "NAME", "UPDATED"}, rows)
		return nil
	},
}

var collectionsInfoCmd = &cobra.Command{
	Use:   "info <id>",
	Short: "Get collection details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		coll, err := client.Collections.Info(getContext(), args[0])
		if err != nil {
			return err
		}
		if outputFormat == "json" {
			printJSON(coll)
		} else {
			fmt.Printf("ID:          %s\n", coll.ID)
			fmt.Printf("Name:        %s\n", coll.Name)
			fmt.Printf("Description: %s\n", coll.Description)
			fmt.Printf("Permission:  %s\n", coll.Permission)
			fmt.Printf("Updated:     %s\n", coll.UpdatedAt.Format("2006-01-02 15:04"))
		}
		return nil
	},
}

var collectionsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a collection",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		name, _ := cmd.Flags().GetString("name")
		desc, _ := cmd.Flags().GetString("description")
		color, _ := cmd.Flags().GetString("color")
		coll, err := client.Collections.Create(getContext(), models.CollectionCreateParams{
			Name:        name,
			Description: desc,
			Color:       color,
		})
		if err != nil {
			return err
		}
		fmt.Printf("Created collection: %s (%s)\n", coll.Name, coll.ID)
		return nil
	},
}

var collectionsUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a collection",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		name, _ := cmd.Flags().GetString("name")
		desc, _ := cmd.Flags().GetString("description")
		color, _ := cmd.Flags().GetString("color")
		coll, err := client.Collections.Update(getContext(), models.CollectionUpdateParams{
			ID:          args[0],
			Name:        name,
			Description: desc,
			Color:       color,
		})
		if err != nil {
			return err
		}
		fmt.Printf("Updated collection: %s\n", coll.Name)
		return nil
	},
}

var collectionsDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a collection",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if !cli.Output.Confirm("Delete this collection and all its documents? This cannot be undone") {
			fmt.Println("Cancelled.")
			return nil
		}
		client, err := getClient()
		if err != nil {
			return err
		}
		if err := client.Collections.Delete(getContext(), args[0]); err != nil {
			return err
		}
		fmt.Println("Collection deleted.")
		return nil
	},
}

var collectionsArchiveCmd = &cobra.Command{
	Use:   "archive <id>",
	Short: "Archive a collection",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		coll, err := client.Collections.Archive(getContext(), args[0])
		if err != nil {
			return err
		}
		fmt.Printf("Archived: %s\n", coll.Name)
		return nil
	},
}

var collectionsRestoreCmd = &cobra.Command{
	Use:   "restore <id>",
	Short: "Restore a collection",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		coll, err := client.Collections.Restore(getContext(), args[0])
		if err != nil {
			return err
		}
		fmt.Printf("Restored: %s\n", coll.Name)
		return nil
	},
}

var collectionsTreeCmd = &cobra.Command{
	Use:   "tree <id>",
	Short: "Show document tree for a collection",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		nodes, err := client.Collections.Documents(getContext(), args[0])
		if err != nil {
			return err
		}
		if outputFormat == "json" {
			printJSON(nodes)
		} else {
			printTree(nodes, 0)
		}
		return nil
	},
}

func printTree(nodes []models.NavigationNode, depth int) {
	for _, n := range nodes {
		prefix := strings.Repeat("  ", depth)
		fmt.Printf("%s%s (%s)\n", prefix, n.Title, n.ID)
		if len(n.Children) > 0 {
			printTree(n.Children, depth+1)
		}
	}
}

func init() {
	collectionsCmd.AddCommand(collectionsListCmd)
	collectionsCmd.AddCommand(collectionsInfoCmd)
	collectionsCmd.AddCommand(collectionsCreateCmd)
	collectionsCmd.AddCommand(collectionsUpdateCmd)
	collectionsCmd.AddCommand(collectionsDeleteCmd)
	collectionsCmd.AddCommand(collectionsArchiveCmd)
	collectionsCmd.AddCommand(collectionsRestoreCmd)
	collectionsCmd.AddCommand(collectionsTreeCmd)

	collectionsCreateCmd.Flags().String("name", "", "Collection name")
	collectionsCreateCmd.Flags().String("description", "", "Collection description")
	collectionsCreateCmd.Flags().String("color", "", "Collection color")
	collectionsUpdateCmd.Flags().String("name", "", "New name")
	collectionsUpdateCmd.Flags().String("description", "", "New description")
	collectionsUpdateCmd.Flags().String("color", "", "New color")
}
