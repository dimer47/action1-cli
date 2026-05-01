package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newOrgCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "org",
		Short: "Manage organizations",
	}

	cmd.AddCommand(
		newOrgListCmd(),
		newOrgCreateCmd(),
		newOrgUpdateCmd(),
		newOrgDeleteCmd(),
	)

	return cmd
}

func newOrgListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List organizations",
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, err := getClient().Get("/organizations", nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}

func newOrgCreateCmd() *cobra.Command {
	var name, description, data string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an organization",
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
			if name != "" {
				body["name"] = name
			}
			if description != "" {
				body["description"] = description
			}
			raw, err := getClient().Post("/organizations", body)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "organization name")
	cmd.Flags().StringVar(&description, "description", "", "description")
	cmd.Flags().StringVar(&data, "data", "", "JSON payload")

	return cmd
}

func newOrgUpdateCmd() *cobra.Command {
	var data string

	cmd := &cobra.Command{
		Use:   "update <orgId>",
		Short: "Update organization settings",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			body, err := parseDataFlag(data)
			if err != nil {
				return err
			}
			raw, err := getClient().Patch("/organizations/"+args[0], body)
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

func newOrgDeleteCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete <orgId>",
		Short: "Delete an organization",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !yes {
				if !confirmAction("delete organization " + args[0]) {
					return nil
				}
			}
			_, err := getClient().Delete("/organizations/" + args[0])
			if err != nil {
				return err
			}
			fmt.Println("Organization deleted.")
			return nil
		},
	}

	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "skip confirmation")

	return cmd
}
