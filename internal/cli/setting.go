package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newSettingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setting",
		Short: "Manage advanced settings",
	}

	cmd.AddCommand(
		newSettingTemplateCmd(),
		newSettingListCmd(),
		newSettingCreateCmd(),
		newSettingGetCmd(),
		newSettingUpdateCmd(),
		newSettingDeleteCmd(),
	)

	return cmd
}

func newSettingTemplateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "template",
		Short: "Browse setting templates",
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "list",
			Short: "List setting templates",
			RunE: func(cmd *cobra.Command, args []string) error {
				raw, err := getClient().Get("/setting-templates/all", nil)
				if err != nil {
					return err
				}
				return printRaw(raw)
			},
		},
		&cobra.Command{
			Use:   "get <templateId>",
			Short: "Get a setting template",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				raw, err := getClient().Get("/setting-templates/all/"+args[0], nil)
				if err != nil {
					return err
				}
				return printRaw(raw)
			},
		},
	)

	return cmd
}

func newSettingListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, err := getClient().Get("/settings/all", nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}

func newSettingCreateCmd() *cobra.Command {
	var data string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new setting",
		RunE: func(cmd *cobra.Command, args []string) error {
			body, err := parseDataFlag(data)
			if err != nil {
				return err
			}
			raw, err := getClient().Post("/settings/all", body)
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

func newSettingGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <settingId>",
		Short: "Get setting configuration",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, err := getClient().Get("/settings/all/"+args[0], nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}

func newSettingUpdateCmd() *cobra.Command {
	var data string

	cmd := &cobra.Command{
		Use:   "update <settingId>",
		Short: "Update a setting",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			body, err := parseDataFlag(data)
			if err != nil {
				return err
			}
			raw, err := getClient().Patch("/settings/all/"+args[0], body)
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

func newSettingDeleteCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete <settingId>",
		Short: "Delete a setting",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !yes {
				if !confirmAction("delete setting " + args[0]) {
					return nil
				}
			}
			_, err := getClient().Delete("/settings/all/" + args[0])
			if err != nil {
				return err
			}
			fmt.Println("Setting deleted.")
			return nil
		},
	}

	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "skip confirmation")

	return cmd
}
