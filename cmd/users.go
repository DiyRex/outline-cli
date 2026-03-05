package cmd

import (
	"fmt"
	"outline-cli/internal/models"

	"github.com/spf13/cobra"
)

var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "Manage users",
}

var usersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List users",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		query, _ := cmd.Flags().GetString("query")
		users, _, err := client.Users.List(getContext(), models.UserListParams{
			Query:            query,
			PaginationParams: models.PaginationParams{Limit: 100},
		})
		if err != nil {
			return err
		}
		var rows [][]string
		for _, u := range users {
			rows = append(rows, []string{u.ID, u.Name, u.Email, u.Role})
		}
		printOutput(users, []string{"ID", "NAME", "EMAIL", "ROLE"}, rows)
		return nil
	},
}

var usersInfoCmd = &cobra.Command{
	Use:   "info [id]",
	Short: "Get user details (defaults to current user)",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		id := ""
		if len(args) > 0 {
			id = args[0]
		}
		user, err := client.Users.Info(getContext(), id)
		if err != nil {
			return err
		}
		if outputFormat == "json" {
			printJSON(user)
		} else {
			fmt.Printf("ID:    %s\n", user.ID)
			fmt.Printf("Name:  %s\n", user.Name)
			fmt.Printf("Email: %s\n", user.Email)
			fmt.Printf("Role:  %s\n", user.Role)
		}
		return nil
	},
}

func init() {
	usersCmd.AddCommand(usersListCmd)
	usersCmd.AddCommand(usersInfoCmd)

	usersListCmd.Flags().String("query", "", "Search query")
}
