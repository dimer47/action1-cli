package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Manage missing updates (patches)",
	}

	cmd.AddCommand(
		newUpdateListCmd(),
		newUpdateGetCmd(),
		newUpdateEndpointsCmd(),
		newUpdateApproveCmd(),
	)

	return cmd
}

func newUpdateListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all missing updates",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			raw, err := getClient().Get("/updates/"+orgID, nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}

func newUpdateGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <packageId>",
		Short: "List updates for a specific package",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			raw, err := getClient().Get(fmt.Sprintf("/updates/%s/%s", orgID, args[0]), nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}

func newUpdateEndpointsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "endpoints <packageId> <versionId>",
		Short: "List endpoints missing a specific update",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			raw, err := getClient().Get(fmt.Sprintf("/updates/%s/%s/versions/%s/endpoints", orgID, args[0], args[1]), nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}

func newUpdateApproveCmd() *cobra.Command {
	var data string

	cmd := &cobra.Command{
		Use:   "approve",
		Short: "Set approval status for updates in bulk",
		Long: `Set approval status for updates in bulk.

Requires the approve_updates permission.

The payload is an array of status changes, each carrying an approval_status
("New", "Approved" or "Declined") and the packages it applies to:

  [
    {
      "approval_status": "Approved",
      "packages": [
        {
          "package_id": "The_Git_Development_Community_Git_1693310149374_builtin",
          "version_id": "2.51.0.2_1759192980389"
        }
      ]
    }
  ]

The field is package_id, singular. The published API reference spells it
packages_id, but the endpoint rejects that with "Property name
[0].packages[0].package_id must be set" — checked against the live API on
2026-09-29.

With --org all, a package scoped to a single organization is updated only for
the organization it belongs to.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			body, err := parseDataFlagAny(data)
			if err != nil {
				return err
			}
			raw, err := getClient().Post(fmt.Sprintf("/updates/%s/approvals", orgID), body)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}

	cmd.Flags().StringVar(&data, "data", "", "JSON payload (inline, @file, or -)")
	_ = cmd.MarkFlagRequired("data")

	return cmd
}
