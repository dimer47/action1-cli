package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newSubscriptionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "subscription",
		Aliases: []string{"sub"},
		Short:   "Manage subscriptions and usage",
	}

	cmd.AddCommand(
		newSubInfoCmd(),
		newSubTrialCmd(),
		newSubQuoteCmd(),
		newSubUsageCmd(),
		newSubUsageOrgsCmd(),
		newSubUsageOrgCmd(),
	)

	return cmd
}

func newSubInfoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "info",
		Short: "Get enterprise license information",
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, err := getClient().Get("/subscription/enterprise", nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}

func newSubTrialCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "trial",
		Short: "Request a free trial or extension",
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, err := getClient().Post("/subscription/enterprise/trial", nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}

func newSubQuoteCmd() *cobra.Command {
	var data string

	cmd := &cobra.Command{
		Use:   "quote",
		Short: "Request a quote",
		RunE: func(cmd *cobra.Command, args []string) error {
			var body interface{}
			if data != "" {
				var err error
				body, err = parseDataFlag(data)
				if err != nil {
					return err
				}
			}
			raw, err := getClient().Post("/subscription/enterprise/quote", body)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}

	cmd.Flags().StringVar(&data, "data", "", "JSON payload")

	return cmd
}

func newSubUsageCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "usage",
		Short: "Get enterprise usage statistics",
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, err := getClient().Get("/subscription/usage/enterprise", nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}

func newSubUsageOrgsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "usage-orgs",
		Short: "Get usage statistics per organization",
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, err := getClient().Get("/subscription/usage/organizations", nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}

func newSubUsageOrgCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "usage-org <orgId>",
		Short: "Get usage statistics for a specific organization",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, err := getClient().Get(fmt.Sprintf("/subscription/usage/organizations/%s", args[0]), nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}
