package cli

import (
	"github.com/spf13/cobra"
)

func newLogCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "log",
		Short: "Manage diagnostic logs",
	}

	cmd.AddCommand(newLogGetCmd())

	return cmd
}

func newLogGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Get diagnostic logs",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			raw, err := getClient().Get("/logs/"+orgID, nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}
