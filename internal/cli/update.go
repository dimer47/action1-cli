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
