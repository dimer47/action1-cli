package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newEnterpriseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enterprise",
		Short: "Manage enterprise settings",
	}

	cmd.AddCommand(
		newEnterpriseGetCmd(),
		newEnterpriseUpdateCmd(),
		newEnterpriseCloseCmd(),
		newEnterpriseRevokeClosureCmd(),
	)

	return cmd
}

func newEnterpriseGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Get enterprise settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, err := getClient().Get("/enterprise", nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}

func newEnterpriseUpdateCmd() *cobra.Command {
	var data string

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update enterprise settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			body, err := parseDataFlag(data)
			if err != nil {
				return err
			}
			raw, err := getClient().Patch("/enterprise", body)
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

func newEnterpriseCloseCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "close",
		Short: "Request enterprise account closure",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !yes {
				if !confirmAction("close the enterprise account (THIS IS IRREVERSIBLE)") {
					return nil
				}
			}
			raw, err := getClient().Post("/enterprise/request-closure", nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}

	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "skip confirmation")

	return cmd
}

func newEnterpriseRevokeClosureCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "revoke-closure",
		Short: "Revoke a pending account closure",
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, err := getClient().Post("/enterprise/revoke-closure", nil)
			if err != nil {
				return err
			}
			fmt.Println("Account closure revoked.")
			return printRaw(raw)
		},
	}
}
