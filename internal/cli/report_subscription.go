package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newReportSubscriptionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "report-subscription",
		Aliases: []string{"report-sub"},
		Short:   "Manage report email subscriptions",
	}

	cmd.AddCommand(
		newReportSubListCmd(),
		newReportSubCreateCmd(),
		newReportSubUpdateCmd(),
		newReportSubDeleteCmd(),
	)

	return cmd
}

func newReportSubListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List report subscriptions",
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, err := getClient().Get("/me/report-subscriptions", nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}

func newReportSubCreateCmd() *cobra.Command {
	var data string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a report subscription",
		RunE: func(cmd *cobra.Command, args []string) error {
			body, err := parseDataFlag(data)
			if err != nil {
				return err
			}
			raw, err := getClient().Post("/me/report-subscriptions", body)
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

func newReportSubUpdateCmd() *cobra.Command {
	var data string

	cmd := &cobra.Command{
		Use:   "update <subscriptionId>",
		Short: "Update a report subscription",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			body, err := parseDataFlag(data)
			if err != nil {
				return err
			}
			raw, err := getClient().Patch("/me/report-subscriptions/"+args[0], body)
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

func newReportSubDeleteCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete <subscriptionId>",
		Short: "Delete a report subscription",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !yes {
				if !confirmAction("delete report subscription " + args[0]) {
					return nil
				}
			}
			_, err := getClient().Delete("/me/report-subscriptions/" + args[0])
			if err != nil {
				return err
			}
			fmt.Println("Report subscription deleted.")
			return nil
		},
	}

	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "skip confirmation")

	return cmd
}
