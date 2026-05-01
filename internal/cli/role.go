package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newRoleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "role",
		Short: "Manage roles (RBAC)",
	}

	cmd.AddCommand(
		newRoleListCmd(),
		newRoleCreateCmd(),
		newRoleGetCmd(),
		newRoleUpdateCmd(),
		newRoleDeleteCmd(),
		newRoleCloneCmd(),
		newRoleUsersCmd(),
		newRoleAssignCmd(),
		newRoleUnassignCmd(),
		newRolePermissionsCmd(),
	)

	return cmd
}

func newRoleListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List roles",
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, err := getClient().Get("/roles", nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}

func newRoleCreateCmd() *cobra.Command {
	var data string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a role",
		RunE: func(cmd *cobra.Command, args []string) error {
			body, err := parseDataFlag(data)
			if err != nil {
				return err
			}
			raw, err := getClient().Post("/roles", body)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}

	cmd.Flags().StringVar(&data, "data", "", "JSON payload")
	_ = cmd.MarkFlagRequired("data")

	return cmd
}

func newRoleGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <roleId>",
		Short: "Get a specific role",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, err := getClient().Get("/roles/"+args[0], nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}

func newRoleUpdateCmd() *cobra.Command {
	var data string

	cmd := &cobra.Command{
		Use:   "update <roleId>",
		Short: "Update a role",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			body, err := parseDataFlag(data)
			if err != nil {
				return err
			}
			raw, err := getClient().Patch("/roles/"+args[0], body)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}

	cmd.Flags().StringVar(&data, "data", "", "JSON payload")
	_ = cmd.MarkFlagRequired("data")

	return cmd
}

func newRoleDeleteCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete <roleId>",
		Short: "Delete a role",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !yes {
				if !confirmAction("delete role " + args[0]) {
					return nil
				}
			}
			_, err := getClient().Delete("/roles/" + args[0])
			if err != nil {
				return err
			}
			fmt.Println("Role deleted.")
			return nil
		},
	}

	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "skip confirmation")

	return cmd
}

func newRoleCloneCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "clone <roleId>",
		Short: "Clone a role",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, err := getClient().Post("/roles/"+args[0]+"/clone", nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}

func newRoleUsersCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "users <roleId>",
		Short: "List users in a role",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, err := getClient().Get("/roles/"+args[0]+"/users", nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}

func newRoleAssignCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "assign <roleId> <userId>",
		Short: "Assign a user to a role",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, err := getClient().Post(fmt.Sprintf("/roles/%s/users/%s", args[0], args[1]), nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}

func newRoleUnassignCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "unassign <roleId> <userId>",
		Short: "Unassign a user from a role",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := getClient().Delete(fmt.Sprintf("/roles/%s/users/%s", args[0], args[1]))
			if err != nil {
				return err
			}
			fmt.Println("User unassigned from role.")
			return nil
		},
	}
}

func newRolePermissionsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "permissions",
		Short: "List permission templates",
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, err := getClient().Get("/permissions", nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}
