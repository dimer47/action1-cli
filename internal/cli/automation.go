package cli

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

func newAutomationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "automation",
		Aliases: []string{"auto"},
		Short:   "Manage automations (schedules, instances, templates)",
	}

	cmd.AddCommand(
		newAutoScheduleCmd(),
		newAutoInstanceCmd(),
		newAutoTemplateCmd(),
	)

	return cmd
}

// --- Schedule subcommands ---

func newAutoScheduleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schedule",
		Short: "Manage automation schedules",
	}

	cmd.AddCommand(
		newAutoScheduleListCmd(),
		newAutoScheduleCreateCmd(),
		newAutoScheduleGetCmd(),
		newAutoScheduleUpdateCmd(),
		newAutoScheduleDeleteCmd(),
		newAutoScheduleDeploymentStatusCmd(),
		newAutoScheduleRemoveActionCmd(),
	)

	return cmd
}

func newAutoScheduleListCmd() *cobra.Command {
	var limit int
	var filter string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List automation schedules",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			q := url.Values{}
			if limit > 0 {
				q.Set("$top", fmt.Sprintf("%d", limit))
			}
			if filter != "" {
				q.Set("$filter", filter)
			}
			raw, err := getClient().Get("/automations/schedules/"+orgID, q)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "l", 0, "max number of results")
	cmd.Flags().StringVarP(&filter, "filter", "f", "", "filter expression")

	return cmd
}

func newAutoScheduleCreateCmd() *cobra.Command {
	var data string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an automation schedule",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			body, err := parseDataFlag(data)
			if err != nil {
				return err
			}
			raw, err := getClient().Post("/automations/schedules/"+orgID, body)
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

func newAutoScheduleGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <automationId>",
		Short: "Get a specific automation schedule",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			raw, err := getClient().Get(fmt.Sprintf("/automations/schedules/%s/%s", orgID, args[0]), nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}

func newAutoScheduleUpdateCmd() *cobra.Command {
	var data string

	cmd := &cobra.Command{
		Use:   "update <automationId>",
		Short: "Update an automation schedule",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			body, err := parseDataFlag(data)
			if err != nil {
				return err
			}
			raw, err := getClient().Patch(fmt.Sprintf("/automations/schedules/%s/%s", orgID, args[0]), body)
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

func newAutoScheduleDeleteCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete <automationId>",
		Short: "Delete an automation schedule",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			if !yes {
				if !confirmAction("delete automation schedule " + args[0]) {
					return nil
				}
			}
			_, err := getClient().Delete(fmt.Sprintf("/automations/schedules/%s/%s", orgID, args[0]))
			if err != nil {
				return err
			}
			fmt.Println("Automation schedule deleted.")
			return nil
		},
	}

	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "skip confirmation")

	return cmd
}

func newAutoScheduleDeploymentStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "deployment-status <automationId>",
		Short: "Get deployment statuses for a schedule",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			raw, err := getClient().Get(fmt.Sprintf("/automations/schedules/%s/%s/deployment-statuses", orgID, args[0]), nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}

func newAutoScheduleRemoveActionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove-action <automationId> <actionId>",
		Short: "Remove an action from a schedule",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			_, err := getClient().Delete(fmt.Sprintf("/automations/schedules/%s/%s/actions/%s", orgID, args[0], args[1]))
			if err != nil {
				return err
			}
			fmt.Println("Action removed.")
			return nil
		},
	}
}

// --- Instance subcommands ---

func newAutoInstanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "instance",
		Short: "Manage automation instances",
	}

	cmd.AddCommand(
		newAutoInstanceListCmd(),
		newAutoInstanceRunCmd(),
		newAutoInstanceGetCmd(),
		newAutoInstanceResultsCmd(),
		newAutoInstanceResultDetailsCmd(),
		newAutoInstanceStopCmd(),
	)

	return cmd
}

func newAutoInstanceListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List automation instances",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			raw, err := getClient().Get("/automations/instances/"+orgID, nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}

func newAutoInstanceRunCmd() *cobra.Command {
	var data string

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Apply (run) an automation immediately",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			body, err := parseDataFlag(data)
			if err != nil {
				return err
			}
			raw, err := getClient().Post("/automations/instances/"+orgID, body)
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

func newAutoInstanceGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <instanceId>",
		Short: "Get automation instance details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			raw, err := getClient().Get(fmt.Sprintf("/automations/instances/%s/%s", orgID, args[0]), nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}

func newAutoInstanceResultsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "results <instanceId>",
		Short: "List endpoint results for an instance",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			raw, err := getClient().Get(fmt.Sprintf("/automations/instances/%s/%s/endpoint-results", orgID, args[0]), nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}

func newAutoInstanceResultDetailsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "result-details <instanceId> <endpointId>",
		Short: "Get detailed results for an endpoint",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			raw, err := getClient().Get(fmt.Sprintf("/automations/instances/%s/%s/endpoint-results/%s/details", orgID, args[0], args[1]), nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}

func newAutoInstanceStopCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stop <instanceId>",
		Short: "Stop a running automation",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			raw, err := getClient().Post(fmt.Sprintf("/automations/instances/%s/%s/stop", orgID, args[0]), nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}

// --- Template subcommands ---

func newAutoTemplateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "template",
		Short: "Browse action templates",
	}

	cmd.AddCommand(
		newAutoTemplateListCmd(),
		newAutoTemplateGetCmd(),
	)

	return cmd
}

func newAutoTemplateListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List action templates",
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, err := getClient().Get("/automations/action-templates", nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}

func newAutoTemplateGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <templateId>",
		Short: "Get a specific action template",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, err := getClient().Get("/automations/action-templates/"+args[0], nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}
