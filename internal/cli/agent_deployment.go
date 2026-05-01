package cli

import (
	"github.com/spf13/cobra"
)

func newAgentDeploymentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agent-deployment",
		Short: "Manage agent deployment settings",
	}

	cmd.AddCommand(
		newAgentDeploymentGetCmd(),
		newAgentDeploymentUpdateCmd(),
	)

	return cmd
}

func newAgentDeploymentGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Get agent deployment settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			raw, err := getClient().Get("/endpoints/agent-deployment/"+orgID, nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}

func newAgentDeploymentUpdateCmd() *cobra.Command {
	var data string

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update agent deployment settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			body, err := parseDataFlag(data)
			if err != nil {
				return err
			}
			raw, err := getClient().Patch("/endpoints/agent-deployment/"+orgID, body)
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
