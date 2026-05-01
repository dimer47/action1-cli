package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newRemoteSessionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remote-session",
		Short: "Manage remote desktop sessions",
	}

	cmd.AddCommand(
		newRemoteSessionStartCmd(),
		newRemoteSessionGetCmd(),
		newRemoteSessionSwitchMonitorCmd(),
	)

	return cmd
}

func newRemoteSessionStartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start <endpointId>",
		Short: "Start a new remote session",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			raw, err := getClient().Post(fmt.Sprintf("/endpoints/managed/%s/%s/remote-sessions", orgID, args[0]), nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}

func newRemoteSessionGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <endpointId> <sessionId>",
		Short: "Get remote session details",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			raw, err := getClient().Get(fmt.Sprintf("/endpoints/managed/%s/%s/remote-sessions/%s", orgID, args[0], args[1]), nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}

func newRemoteSessionSwitchMonitorCmd() *cobra.Command {
	var data string

	cmd := &cobra.Command{
		Use:   "switch-monitor <endpointId> <sessionId>",
		Short: "Switch current monitor in a remote session",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			body, err := parseDataFlag(data)
			if err != nil {
				return err
			}
			raw, err := getClient().Patch(fmt.Sprintf("/endpoints/managed/%s/%s/remote-sessions/%s", orgID, args[0], args[1]), body)
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
