package cli

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

func newEndpointCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "endpoint",
		Aliases: []string{"ep"},
		Short:   "Manage endpoints",
	}

	cmd.AddCommand(
		newEndpointListCmd(),
		newEndpointGetCmd(),
		newEndpointStatusCmd(),
		newEndpointUpdateCmd(),
		newEndpointDeleteCmd(),
		newEndpointMoveCmd(),
		newEndpointMissingUpdatesCmd(),
		newEndpointInstallURLCmd(),
	)

	return cmd
}

func newEndpointListCmd() *cobra.Command {
	var limit int
	var filter string
	var all bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all endpoints",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}

			client := getClient()
			q := url.Values{}
			if limit > 0 {
				q.Set("$top", fmt.Sprintf("%d", limit))
			}
			if filter != "" {
				q.Set("$filter", filter)
			}

			if all {
				items, err := client.GetAll("/endpoints/managed/"+orgID, q)
				if err != nil {
					return err
				}
				return printResult(rawToInterface(items))
			}

			raw, err := client.Get("/endpoints/managed/"+orgID, q)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "l", 0, "max number of results")
	cmd.Flags().StringVarP(&filter, "filter", "f", "", "OData filter expression")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all results (auto-paginate)")

	return cmd
}

func newEndpointGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <endpointId>",
		Short: "Get a specific endpoint",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			raw, err := getClient().Get(fmt.Sprintf("/endpoints/managed/%s/%s", orgID, args[0]), nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}

func newEndpointStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Check endpoint status",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			raw, err := getClient().Get("/endpoints/status/"+orgID, nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}

func newEndpointUpdateCmd() *cobra.Command {
	var name, comment string

	cmd := &cobra.Command{
		Use:   "update <endpointId>",
		Short: "Update endpoint name or comment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}

			body := map[string]interface{}{}
			if name != "" {
				body["name"] = name
			}
			if comment != "" {
				body["comment"] = comment
			}

			raw, err := getClient().Patch(fmt.Sprintf("/endpoints/managed/%s/%s", orgID, args[0]), body)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "new name")
	cmd.Flags().StringVar(&comment, "comment", "", "new comment")

	return cmd
}

func newEndpointDeleteCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete <endpointId>",
		Short: "Delete an endpoint",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			if !yes {
				if !confirmAction("delete endpoint " + args[0]) {
					return nil
				}
			}
			_, err := getClient().Delete(fmt.Sprintf("/endpoints/managed/%s/%s", orgID, args[0]))
			if err != nil {
				return err
			}
			fmt.Println("Endpoint deleted.")
			return nil
		},
	}

	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "skip confirmation")

	return cmd
}

func newEndpointMoveCmd() *cobra.Command {
	var toOrg string

	cmd := &cobra.Command{
		Use:   "move <endpointId>",
		Short: "Move endpoint to another organization",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			if toOrg == "" {
				return fmt.Errorf("--to-org is required")
			}

			body := map[string]interface{}{
				"organization_id": toOrg,
			}
			raw, err := getClient().Post(fmt.Sprintf("/endpoints/managed/%s/%s/move", orgID, args[0]), body)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}

	cmd.Flags().StringVar(&toOrg, "to-org", "", "target organization ID")
	_ = cmd.MarkFlagRequired("to-org")

	return cmd
}

func newEndpointMissingUpdatesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "missing-updates <endpointId>",
		Short: "List missing updates for an endpoint",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			raw, err := getClient().Get(fmt.Sprintf("/endpoints/managed/%s/%s/missing-updates", orgID, args[0]), nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}

func newEndpointInstallURLCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "install-url <type>",
		Short: "Get agent installation URL (windowsEXE|windowsMSI|macOS|linux)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			raw, err := getClient().Get(fmt.Sprintf("/endpoints/agent-installation/%s/%s", orgID, args[0]), nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}
