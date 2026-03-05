package cmd

import (
	"fmt"
	"outline-cli/internal/models"

	"github.com/spf13/cobra"
)

var groupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "Manage groups",
}

var groupsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List groups",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		groups, _, err := client.Groups.List(getContext(), models.GroupListParams{
			PaginationParams: models.PaginationParams{Limit: 100},
		})
		if err != nil {
			return err
		}
		var rows [][]string
		for _, g := range groups {
			rows = append(rows, []string{g.ID, g.Name, fmt.Sprintf("%d", g.MemberCount)})
		}
		printOutput(groups, []string{"ID", "NAME", "MEMBERS"}, rows)
		return nil
	},
}

var groupsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a group",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		name, _ := cmd.Flags().GetString("name")
		group, err := client.Groups.Create(getContext(), models.GroupCreateParams{Name: name})
		if err != nil {
			return err
		}
		fmt.Printf("Created group: %s (%s)\n", group.Name, group.ID)
		return nil
	},
}

var groupsDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a group",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		if err := client.Groups.Delete(getContext(), args[0]); err != nil {
			return err
		}
		fmt.Println("Group deleted.")
		return nil
	},
}

var groupsMembersCmd = &cobra.Command{
	Use:   "members <id>",
	Short: "List group members",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		members, _, err := client.Groups.Members(getContext(), args[0], models.PaginationParams{Limit: 100})
		if err != nil {
			return err
		}
		var rows [][]string
		for _, m := range members {
			rows = append(rows, []string{m.User.ID, m.User.Name, m.User.Email})
		}
		printOutput(members, []string{"USER ID", "NAME", "EMAIL"}, rows)
		return nil
	},
}

var groupsAddUserCmd = &cobra.Command{
	Use:   "add-user <group-id> <user-id>",
	Short: "Add a user to a group",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		if err := client.Groups.AddUser(getContext(), args[0], args[1]); err != nil {
			return err
		}
		fmt.Println("User added to group.")
		return nil
	},
}

var groupsRemoveUserCmd = &cobra.Command{
	Use:   "remove-user <group-id> <user-id>",
	Short: "Remove a user from a group",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}
		if err := client.Groups.RemoveUser(getContext(), args[0], args[1]); err != nil {
			return err
		}
		fmt.Println("User removed from group.")
		return nil
	},
}

func init() {
	groupsCmd.AddCommand(groupsListCmd)
	groupsCmd.AddCommand(groupsCreateCmd)
	groupsCmd.AddCommand(groupsDeleteCmd)
	groupsCmd.AddCommand(groupsMembersCmd)
	groupsCmd.AddCommand(groupsAddUserCmd)
	groupsCmd.AddCommand(groupsRemoveUserCmd)

	groupsCreateCmd.Flags().String("name", "", "Group name")
}
