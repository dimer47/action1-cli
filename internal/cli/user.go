package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Manage users",
	}

	cmd.AddCommand(
		newUserMeCmd(),
		newUserMeUpdateCmd(),
		newUserListCmd(),
		newUserCreateCmd(),
		newUserGetCmd(),
		newUserUpdateCmd(),
		newUserDeleteCmd(),
		newUserRolesCmd(),
	)

	return cmd
}

func newUserMeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "me",
		Short: "Get current user",
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, err := getClient().Get("/me", nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}

func newUserMeUpdateCmd() *cobra.Command {
	var data string

	cmd := &cobra.Command{
		Use:   "me-update",
		Short: "Update current user settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			body, err := parseDataFlag(data)
			if err != nil {
				return err
			}
			raw, err := getClient().Patch("/me", body)
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

func newUserListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List users",
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, err := getClient().Get("/users", nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}

func newUserCreateCmd() *cobra.Command {
	var email, name, data string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new user",
		RunE: func(cmd *cobra.Command, args []string) error {
			var body map[string]interface{}
			if data != "" {
				var err error
				body, err = parseDataFlag(data)
				if err != nil {
					return err
				}
			} else {
				body = map[string]interface{}{}
			}
			if email != "" {
				body["email"] = email
			}
			if name != "" {
				body["name"] = name
			}
			raw, err := getClient().Post("/users", body)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}

	cmd.Flags().StringVar(&email, "email", "", "user email")
	cmd.Flags().StringVar(&name, "name", "", "user name")
	cmd.Flags().StringVar(&data, "data", "", "JSON payload")

	return cmd
}

func newUserGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <userId>",
		Short: "Get user details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, err := getClient().Get("/users/"+args[0], nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}

func newUserUpdateCmd() *cobra.Command {
	var data string

	cmd := &cobra.Command{
		Use:   "update <userId>",
		Short: "Update a user",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			body, err := parseDataFlag(data)
			if err != nil {
				return err
			}
			raw, err := getClient().Patch("/users/"+args[0], body)
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

func newUserDeleteCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete <userId>",
		Short: "Delete a user",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !yes {
				if !confirmAction("delete user " + args[0]) {
					return nil
				}
			}
			_, err := getClient().Delete("/users/" + args[0])
			if err != nil {
				return err
			}
			fmt.Println("User deleted.")
			return nil
		},
	}

	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "skip confirmation")

	return cmd
}

func newUserRolesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "roles <userId>",
		Short: "List roles assigned to a user",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, err := getClient().Get("/users/"+args[0]+"/roles", nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}
