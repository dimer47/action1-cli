package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newDataSourceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "data-source",
		Aliases: []string{"ds"},
		Short:   "Manage data sources",
	}

	cmd.AddCommand(
		newDSListCmd(),
		newDSCreateCmd(),
		newDSGetCmd(),
		newDSUpdateCmd(),
		newDSDeleteCmd(),
	)

	return cmd
}

func newDSListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List data sources",
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, err := getClient().Get("/data-sources/all", nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}

func newDSCreateCmd() *cobra.Command {
	var data string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a data source",
		RunE: func(cmd *cobra.Command, args []string) error {
			body, err := parseDataFlag(data)
			if err != nil {
				return err
			}
			raw, err := getClient().Post("/data-sources/all", body)
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

func newDSGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <dataSourceId>",
		Short: "Get a specific data source",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, err := getClient().Get("/data-sources/all/"+args[0], nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}

func newDSUpdateCmd() *cobra.Command {
	var data string

	cmd := &cobra.Command{
		Use:   "update <dataSourceId>",
		Short: "Update a custom data source",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			body, err := parseDataFlag(data)
			if err != nil {
				return err
			}
			raw, err := getClient().Patch("/data-sources/all/"+args[0], body)
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

func newDSDeleteCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete <dataSourceId>",
		Short: "Delete a custom data source",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !yes {
				if !confirmAction("delete data source " + args[0]) {
					return nil
				}
			}
			_, err := getClient().Delete("/data-sources/all/" + args[0])
			if err != nil {
				return err
			}
			fmt.Println("Data source deleted.")
			return nil
		},
	}

	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "skip confirmation")

	return cmd
}
