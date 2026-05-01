package cli

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/spf13/cobra"
)

func newEndpointGroupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "endpoint-group",
		Aliases: []string{"epg"},
		Short:   "Manage endpoint groups",
	}

	cmd.AddCommand(
		newEPGListCmd(),
		newEPGCreateCmd(),
		newEPGGetCmd(),
		newEPGUpdateCmd(),
		newEPGDeleteCmd(),
		newEPGMembersCmd(),
		newEPGAddCmd(),
		newEPGRemoveCmd(),
	)

	return cmd
}

func newEPGListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List endpoint groups",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			raw, err := getClient().Get("/endpoints/groups/"+orgID, nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}

func newEPGCreateCmd() *cobra.Command {
	var name, data string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an endpoint group",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}

			var body map[string]interface{}
			if data != "" {
				var err error
				body, err = parseDataFlag(data)
				if err != nil {
					return err
				}
			} else {
				body = map[string]interface{}{}
			}

			if name != "" {
				body["name"] = name
			}

			raw, err := getClient().Post("/endpoints/groups/"+orgID, body)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "group name")
	cmd.Flags().StringVar(&data, "data", "", "JSON payload (inline, @file, or -)")

	return cmd
}

func newEPGGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <groupId>",
		Short: "Get an endpoint group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			raw, err := getClient().Get(fmt.Sprintf("/endpoints/groups/%s/%s", orgID, args[0]), nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}

func newEPGUpdateCmd() *cobra.Command {
	var data string

	cmd := &cobra.Command{
		Use:   "update <groupId>",
		Short: "Update endpoint group settings",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			body, err := parseDataFlag(data)
			if err != nil {
				return err
			}
			raw, err := getClient().Patch(fmt.Sprintf("/endpoints/groups/%s/%s", orgID, args[0]), body)
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

func newEPGDeleteCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete <groupId>",
		Short: "Delete an endpoint group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			if !yes {
				if !confirmAction("delete endpoint group " + args[0]) {
					return nil
				}
			}
			_, err := getClient().Delete(fmt.Sprintf("/endpoints/groups/%s/%s", orgID, args[0]))
			if err != nil {
				return err
			}
			fmt.Println("Endpoint group deleted.")
			return nil
		},
	}

	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "skip confirmation")

	return cmd
}

func newEPGMembersCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "members <groupId>",
		Short: "List endpoints in a group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			raw, err := getClient().Get(fmt.Sprintf("/endpoints/groups/%s/%s/contents", orgID, args[0]), nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}

func newEPGAddCmd() *cobra.Command {
	var endpoints string

	cmd := &cobra.Command{
		Use:   "add <groupId>",
		Short: "Add endpoints to a group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			ids := strings.Split(endpoints, ",")
			body := map[string]interface{}{
				"add": ids,
			}
			raw, err := getClient().Post(fmt.Sprintf("/endpoints/groups/%s/%s/contents", orgID, args[0]), body)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}

	cmd.Flags().StringVar(&endpoints, "endpoints", "", "comma-separated endpoint IDs")
	_ = cmd.MarkFlagRequired("endpoints")

	return cmd
}

func newEPGRemoveCmd() *cobra.Command {
	var endpoints string

	cmd := &cobra.Command{
		Use:   "remove <groupId>",
		Short: "Remove endpoints from a group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}

			q := url.Values{}
			ids := strings.Split(endpoints, ",")
			body := map[string]interface{}{
				"remove": ids,
			}
			raw, err := getClient().Post(fmt.Sprintf("/endpoints/groups/%s/%s/contents", orgID, args[0]), body)
			_ = q
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}

	cmd.Flags().StringVar(&endpoints, "endpoints", "", "comma-separated endpoint IDs")
	_ = cmd.MarkFlagRequired("endpoints")

	return cmd
}
