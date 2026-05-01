package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage CLI configuration",
	}

	cmd.AddCommand(
		newConfigInitCmd(),
		newConfigShowCmd(),
		newConfigSetCmd(),
		newConfigGetCmd(),
		newConfigUnsetCmd(),
		newConfigListProfilesCmd(),
		newConfigUseProfileCmd(),
	)

	return cmd
}

func newConfigInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize a new configuration file",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cfg.Save(cfgFile); err != nil {
				return err
			}
			fmt.Printf("Configuration saved to %s\n", cfgFile)
			return nil
		},
	}
}

func newConfigShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Show current configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Config file: %s\n", cfgFile)
			fmt.Printf("Active profile: %s\n", cfg.CurrentProfile)
			fmt.Println()

			p := cfg.ActiveProfile()
			fmt.Printf("Region:  %s\n", p.Region)
			fmt.Printf("Org:     %s\n", p.Org)
			fmt.Printf("Output:  %s\n", p.Output)
			return nil
		},
	}
}

func newConfigSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a configuration value",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cfg.SetProfileValue(cfg.CurrentProfile, args[0], args[1]); err != nil {
				return err
			}
			if err := cfg.Save(cfgFile); err != nil {
				return err
			}
			fmt.Printf("Set %s = %s\n", args[0], args[1])
			return nil
		},
	}
}

func newConfigGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <key>",
		Short: "Get a configuration value",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			val, err := cfg.GetProfileValue(cfg.CurrentProfile, args[0])
			if err != nil {
				return err
			}
			fmt.Println(val)
			return nil
		},
	}
}

func newConfigUnsetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "unset <key>",
		Short: "Remove a configuration value",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cfg.SetProfileValue(cfg.CurrentProfile, args[0], ""); err != nil {
				return err
			}
			if err := cfg.Save(cfgFile); err != nil {
				return err
			}
			fmt.Printf("Unset %s\n", args[0])
			return nil
		},
	}
}

func newConfigListProfilesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list-profiles",
		Short: "List available profiles",
		RunE: func(cmd *cobra.Command, args []string) error {
			for name := range cfg.Profiles {
				marker := "  "
				if name == cfg.CurrentProfile {
					marker = "* "
				}
				fmt.Printf("%s%s\n", marker, name)
			}
			return nil
		},
	}
}

func newConfigUseProfileCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "use-profile <name>",
		Short: "Switch to a different profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			if _, ok := cfg.Profiles[name]; !ok {
				// Create the profile if it doesn't exist
				cfg.Profiles[name] = cfg.Profiles[cfg.CurrentProfile]
				fmt.Printf("Created new profile %q\n", name)
			}
			cfg.CurrentProfile = name
			if err := cfg.Save(cfgFile); err != nil {
				return err
			}
			fmt.Printf("Switched to profile %q\n", name)
			return nil
		},
	}
}
