package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newDeployerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deployer",
		Short: "Manage deployers",
	}

	cmd.AddCommand(
		newDeployerListCmd(),
		newDeployerGetCmd(),
		newDeployerDeleteCmd(),
		newDeployerInstallURLCmd(),
	)

	return cmd
}

func newDeployerListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List deployers",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			raw, err := getClient().Get("/endpoints/deployers/"+orgID, nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}

func newDeployerGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <deployerId>",
		Short: "Get a specific deployer",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			raw, err := getClient().Get(fmt.Sprintf("/endpoints/deployers/%s/%s", orgID, args[0]), nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}

func newDeployerDeleteCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete <deployerId>",
		Short: "Delete a deployer",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			if !yes {
				if !confirmAction("delete deployer " + args[0]) {
					return nil
				}
			}
			_, err := getClient().Delete(fmt.Sprintf("/endpoints/deployers/%s/%s", orgID, args[0]))
			if err != nil {
				return err
			}
			fmt.Println("Deployer deleted.")
			return nil
		},
	}

	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "skip confirmation")

	return cmd
}

func newDeployerInstallURLCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "install-url",
		Short: "Get deployer installation URL (Windows EXE)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			raw, err := getClient().Get(fmt.Sprintf("/endpoints/deployer-installation/%s/windowsEXE", orgID), nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}
