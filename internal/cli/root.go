package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/dimer47/action1-cli/internal/api"
	"github.com/dimer47/action1-cli/internal/auth"
	"github.com/dimer47/action1-cli/internal/config"
	"github.com/dimer47/action1-cli/internal/output"
)

var (
	cfgFile    string
	profile    string
	orgID      string
	region     string
	outputFmt  string
	quiet      bool
	verbose    bool
	noColor    bool
	noKeychain bool

	cfg       *config.Config
	apiClient *api.Client
	store     auth.Store
)

// Version is set at build time.
var Version = "dev"

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "action1",
		Short: "CLI for the Action1 platform",
		Long:  "Command-line interface for managing endpoints, automations, reports, and more via the Action1 API.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Skip config loading for completion commands
			if cmd.Name() == "completion" || cmd.Name() == "__complete" {
				return nil
			}

			var err error
			cfg, err = config.Load(cfgFile)
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			if profile == "" {
				profile = cfg.CurrentProfile
			}

			p := cfg.ActiveProfile()

			if region == "" {
				region = string(p.Region)
			}
			if orgID == "" {
				orgID = p.Org
			}
			if outputFmt == "" {
				outputFmt = p.Output
			}
			if outputFmt == "" {
				outputFmt = "table"
			}

			store = auth.NewStore(cfgFile, noKeychain)

			return nil
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file path (default: auto-detected)")
	rootCmd.PersistentFlags().StringVarP(&profile, "profile", "p", "", "config profile to use (default: from config)")
	rootCmd.PersistentFlags().StringVarP(&orgID, "org", "o", "", "organization ID (default: from config)")
	rootCmd.PersistentFlags().StringVarP(&region, "region", "r", "", "server region: na|eu|au (default: from config)")
	rootCmd.PersistentFlags().StringVarP(&outputFmt, "output", "O", "", "output format: table|json|csv|yaml (default: from config)")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "suppress headers and decorations")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "show HTTP requests")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "disable colored output")
	rootCmd.PersistentFlags().BoolVar(&noKeychain, "no-keychain", false, "use file-based credential storage")

	// Register all subcommands
	rootCmd.AddCommand(
		newAuthCmd(),
		newConfigCmd(),
		newEndpointCmd(),
		newEndpointGroupCmd(),
		newRemoteSessionCmd(),
		newDeployerCmd(),
		newAgentDeploymentCmd(),
		newAutomationCmd(),
		newReportCmd(),
		newReportSubscriptionCmd(),
		newSoftwareCmd(),
		newUpdateCmd(),
		newInstalledSoftwareCmd(),
		newVulnerabilityCmd(),
		newDataSourceCmd(),
		newScriptCmd(),
		newSettingCmd(),
		newOrgCmd(),
		newUserCmd(),
		newRoleCmd(),
		newEnterpriseCmd(),
		newSubscriptionCmd(),
		newSearchCmd(),
		newLogCmd(),
		newAuditCmd(),
		newVersionCmd(),
	)

	return rootCmd
}

// Execute runs the root command.
func Execute() {
	cmd := NewRootCmd()
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// Helper: get or create the API client.
func getClient() *api.Client {
	if apiClient == nil {
		apiClient = api.NewClient(config.Region(region), store, profile, verbose)
	}
	return apiClient
}

// Helper: require --org to be set.
func requireOrg() error {
	if orgID == "" {
		return fmt.Errorf("--org is required (or set a default org with 'action1 config set org <id>')")
	}
	return nil
}

// Helper: output formatting.
func printResult(data interface{}) error {
	return output.Print(os.Stdout, output.ParseFormat(outputFmt), data, quiet)
}

func printRaw(raw []byte) error {
	return output.PrintRaw(os.Stdout, output.ParseFormat(outputFmt), raw, quiet)
}
