package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newInstalledSoftwareCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "installed-software",
		Aliases: []string{"isw"},
		Short:   "Manage installed software inventory",
	}

	cmd.AddCommand(
		newISWListCmd(),
		newISWGetCmd(),
		newISWErrorsCmd(),
		newISWRequeryCmd(),
	)

	return cmd
}

func newISWListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List installed apps",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			raw, err := getClient().Get(fmt.Sprintf("/installed-software/%s/data", orgID), nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}

func newISWGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <endpointId>",
		Short: "List installed apps on a specific endpoint",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			raw, err := getClient().Get(fmt.Sprintf("/installed-software/%s/data/%s", orgID, args[0]), nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}

func newISWErrorsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "errors",
		Short: "Get collection errors",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			raw, err := getClient().Get(fmt.Sprintf("/installed-software/%s/errors", orgID), nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}

func newISWRequeryCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "requery [endpointId]",
		Short: "Re-query installed apps (all or specific endpoint)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			path := fmt.Sprintf("/installed-software/%s/requery", orgID)
			if len(args) > 0 {
				path = fmt.Sprintf("/installed-software/%s/requery/%s", orgID, args[0])
			}
			raw, err := getClient().Post(path, nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}
